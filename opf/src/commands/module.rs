use super::{Command, CommandAction, CommandObject};
use crate::common::search::Target as TargetSearch;
use crate::common::Link;
use crate::common::Target;
use crate::compiled::Compiled;
use crate::core::session::Session;
use crate::error::{ErrorKind, Module as ModuleError};
use crate::modules::metadata::{Arg, Args};
use crate::modules::{Manager, Module, OpfModule};
use crate::vm;
use std::collections::HashMap;
use std::ops::Add;

pub fn exec<'a>(
    session: &mut Session,
    cmd: Command,
    manager: &mut Manager,
) -> Result<(), ErrorKind> {
    match cmd.action {
        CommandAction::Run => run(session, cmd, manager),
        CommandAction::List => list(session, manager),
        CommandAction::Help => help(session, cmd, manager),
        _ => Ok(()),
    }
}

pub fn list<'a>(session: &'a Session, manager: &'a Manager) -> Result<(), ErrorKind> {
    let headers = vec![
        "name".to_string(),
        "author".to_string(),
        "description".to_string(),
    ];
    let mut rows = vec![];

    for (_, module) in &manager.modules {
        let fields = vec![module.name(), module.author(), module.resume()];
        rows.push(fields);
    }

    session.output_table(headers, rows);
    Ok(())
}

pub fn help<'a>(
    session: &'a mut Session,
    cmd: Command,
    manager: &'a Manager,
) -> Result<(), ErrorKind> {
    let module_name = match cmd.object {
        CommandObject::Module(ref module_name) => module_name.clone(),
        _ => return Err(ErrorKind::Module(ModuleError::ModuleNameNotFound)),
    };

    for (_, module) in &manager.modules {
        if module_name.eq(&module.name()) {
            let headers = vec![
                "name".to_string(),
                "author".to_string(),
                "description".to_string(),
            ];
            let mut rows = vec![];

            let fields = vec![module.name(), module.author(), module.resume()];
            rows.push(fields);
            session.output_table(headers, rows);

            let headers = vec![
                "Argument".to_string(),
                "default".to_string(),
                "optional".to_string(),
            ];

            let rows = module
                .args()
                .iter()
                .map(|arg| {
                    vec![
                        arg.name.clone(),
                        arg.value.clone().unwrap_or("<None>".to_string()),
                        format!("{}", arg.is_optional),
                    ]
                })
                .collect::<Vec<Vec<String>>>();
            session.output_table(headers, rows);
        }
    }

    Ok(())
}

pub fn run<'a>(
    session: &'a mut Session,
    cmd: Command,
    manager: &'a Manager,
) -> Result<(), ErrorKind> {
    let module_name = match cmd.object {
        CommandObject::Module(module_name) => module_name,
        _ => return Err(ErrorKind::Module(ModuleError::ModuleNameNotFound)),
    };

    if let Some(module) = manager.modules.get(module_name.as_str()) {
        if let Some(group_id) = cmd.params.get("group_id") {
            if !session.exist_group_id(group_id) {
                return Err(ErrorKind::Module(ModuleError::GroupNotFound));
            }
            let get_group = session.get_group_by_id(group_id.as_str());
            if get_group.is_err() {
                return Err(ErrorKind::Module(ModuleError::GroupNotFound));
            }

            let group = get_group.unwrap();
            let targets = session.get_targets_by_group(&group);

            for target in targets {
                let mut new_param = cmd.params.clone();
                new_param.insert("target".to_string(), target.target_id.to_string());
                execute_module(session, &module, &new_param)?;
            }
            return Ok(());
        }
        execute_module(session, &module, &cmd.params)?;
        return Ok(());
    }
    Err(ErrorKind::Module(ModuleError::ModuleNameNotFound))
}

fn verify_arguments(
    session: &Session,
    params: &mut Vec<Arg>,
    args: &HashMap<String, String>,
) -> Result<Target, ErrorKind> {
    let mut target: Option<Target> = None;

    for mut param in params.into_iter() {
        match args.get(&param.name) {
            Some(value) => {
                if param.is_target {
                    let mut search = TargetSearch::default();
                    search.target_id = Some(value.clone());

                    if !session.exist_target(&search) {
                        search.target_id = None;
                        search.target_custom_id = Some(value.clone());

                        if session.exist_target(&search) {
                            break;
                        }

                        return Err(ErrorKind::Module(ModuleError::TargetNotAvailable));
                    }

                    let targets = session.get_targets(search);
                    if targets.len() < 1 {
                        return Err(ErrorKind::Module(ModuleError::TargetNotAvailable));
                    }

                    target = Some(targets[0].clone());
                    param.value = Some(targets[0].target_id.clone().to_string());
                    continue;
                }

                param.value = Some(value.clone());
            }
            None => {
                if !param.is_optional {
                    return Err(ErrorKind::Module(ModuleError::ParamNotAvailable(
                        param.name.clone(),
                    )));
                }
            }
        }
    }

    if target.is_none() {
        return Err(ErrorKind::Module(ModuleError::ParamNotAvailable(
            "target specifier".to_string(),
        )));
    }

    Ok(target.unwrap())
}

fn execute_lua(
    sess: &mut Session,
    module: &Module,
    mut params: Vec<Arg>,
) -> Result<Vec<Target>, ErrorKind> {
    let mut targets: Vec<HashMap<String, String>> = vec![];
    let vm = rlua::Lua::new();

    let _ = vm.context(|mut vm_context| {
        let _ = vm_context.globals().set("targets", targets.clone());

        let arguments = vm::arg::Args::from(params);
        vm_context.globals().set("args", arguments).map_err(|_| {
            ErrorKind::Module(ModuleError::Execution("can't retrieve args".to_string()))
        })?;
        vm_context
            .globals()
            .set("sess_targets", sess.targets.clone())
            .map_err(|_| {
                ErrorKind::Module(ModuleError::Execution(
                    "can't retrieve session targets".to_string(),
                ))
            })?;

        for extend in &module.metadata.extends {
            vm::extends(&mut vm_context, extend, sess);
        }

        vm::require::common::extends_common(&mut vm_context);

        let content = std::fs::read(&module.file_name)
            .map_err(|_| ErrorKind::Module(ModuleError::CantLoadContent))?;

        _ = vm_context.load(content.as_slice()).exec();

        vm_context.scope(|_scope| {
            let new_targets: Vec<HashMap<String, String>> =
                vm_context.globals().get("targets").unwrap_or(vec![]);

            targets = new_targets;
        });
        Ok(())
    })?;

    let mut results = vec![];
    for target in targets {
        let mut new_target = Target::try_from(target).map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(
                "missing argument in target creation".to_string(),
            ))
        })?;

        results.push(new_target);
    }
    Ok(results)
}

fn execute_module(
    sess: &mut Session,
    module: &OpfModule,
    args: &HashMap<String, String>,
) -> Result<(), ErrorKind> {
    let mut params = module.args();
    let target_parent = verify_arguments(sess, &mut params, args)?;

    let targets  = match module {
        OpfModule::Lua(ref lua_module) => execute_lua(sess, lua_module, module.args())?,
        OpfModule::Compiled(ref module) => module.run(sess, Args::from(params))?,
    };

    let mut inserted = 0;
    let module_name = module.name();
    for mut target in targets {

        let group_name = String::from(&target_parent.target_name)
            .add(".")
            .add(&module_name);

        let group = match sess.create_group(&group_name) {
            Ok(group) => group,
            Err(_) => match sess.get_group(&group_name) {
                Ok(group) => group,
                Err(_) => continue,
            },
        };

        target.target_parent = Some(target_parent.target_uuid);
        target.target_groups.push(group.group_uuid);

        if let Ok(target) = sess.create_target(target) {
            if let Some(parent) = target.target_parent {
                let _ = Link::new(parent, target.target_uuid, &module_name, &module_name);
            }
            inserted += 1;
        }
    }

    println!("INF   new targets {}", inserted);

    Ok(())
}

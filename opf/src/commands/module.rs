use super::{Command, CommandAction, CommandObject};
use crate::common::module as opf_module;
use crate::common::search::Target as TargetSearch;
use crate::common::Target;
use crate::core::session::Session;
use crate::modules::metadata::Arg;
use crate::modules::{Manager, Module};
use crate::vm;
use std::collections::HashMap;
use std::ops::Add;

pub fn exec<'a>(
    session: &'a mut Session,
    cmd: Command,
    manager: &'a mut Manager,
) -> Result<(), opf_module::Error> {
    match cmd.action {
        CommandAction::Run => run(session, cmd, manager),
        CommandAction::List => list(session, manager),
        _ => Ok(()),
    }
}

pub fn list<'a>(session: &'a Session, manager: &'a Manager) -> Result<(), opf_module::Error> {
    let headers = vec![
        "name".to_string(),
        "author".to_string(),
        "description".to_string(),
    ];
    let mut rows = vec![];

    for module in &manager.modules {
        let fields = vec![
            module.metadata.name.clone().unwrap_or("-".to_string()),
            module.metadata.author.clone().unwrap_or("-".to_string()),
            module
                .metadata
                .description
                .clone()
                .unwrap_or("-".to_string()),
        ];
        rows.push(fields);
    }

    session.output_table(headers, rows);
    Ok(())
}

pub fn run<'a>(
    session: &'a mut Session,
    cmd: Command,
    manager: &'a Manager,
) -> Result<(), opf_module::Error> {
    let module_name = match cmd.object {
        CommandObject::Module(ref module_name) => module_name.clone(),
        _ => return Err(opf_module::Error::ModuleNameNotFound),
    };

    for module in &manager.modules {
        if let Some(ref name) = module.metadata.name {
            if name.eq(&module_name) {
                if let Some(group_id) = cmd.params.get("group_id") {
                    if !session.exist_group_id(group_id) {
                        return Err(opf_module::Error::GroupNotFound);
                    }
                    let get_group = session.get_group_by_id(group_id.as_str());
                    if get_group.is_err() {
                        return Err(opf_module::Error::GroupNotFound);
                    }

                    let group = get_group.unwrap();
                    let targets = session.get_targets_by_group(&group);

                    for target in targets {
                        let mut new_param = cmd.params.clone();
                        new_param.insert("target".to_string(), target.target_id.to_string());
                        execute_module(session, &module, new_param);
                    }
                    return Ok(());
                }
                execute_module(session, &module, cmd.params.clone());
                return Ok(());
            }
        }
    }

    Err(opf_module::Error::ModuleNameNotFound)
}

fn verify_arguments<'a, 'b, 'c>(
    session: &'a Session,
    params: &'b mut Vec<Arg>,
    args: &'c HashMap<String, String>,
) -> Result<Target, opf_module::Error> {
    let mut target: Option<Target> = None;

    for mut param in params.into_iter() {
        match args.get(&param.name) {
            Some(value) => {
                if param.is_target {
                    let mut search = TargetSearch::default();
                    search.target_id = Some(value.clone());

                    if !session.exist_target(&search) {
                        println!("not exist with id, searching with custom_id...");
                        search.target_id = None;
                        search.target_custom_id = Some(value.clone());

                        if session.exist_target(&search) {
                            println!("found with custom id = {:?}", search);
                            break;
                        }

                        return Err(opf_module::Error::TargetNotAvailable);
                    }

                    let targets = session.get_targets(search);
                    if targets.len() < 1 {
                        return Err(opf_module::Error::TargetNotAvailable);
                    }

                    target = Some(targets[0].clone());
                    param.value = Some(targets[0].target_id.clone().to_string());
                    continue;
                }

                param.value = Some(value.clone());
            }
            None => {
                if !param.is_optional {
                    return Err(opf_module::Error::ParamNotAvailable(param.name.clone()));
                }
            }
        }
    }

    if target.is_none() {
        return Err(opf_module::Error::ParamNotAvailable(
            "target specifier".to_string(),
        ));
    }

    Ok(target.unwrap())
}

fn execute_module<'a, 'b>(
    sess: &'a mut Session,
    module: &'b Module,
    args: HashMap<String, String>,
) {
    let mut targets: Vec<HashMap<String, String>> = vec![];

    // checking arguments
    let mut params = module.metadata.args.clone();

    let target_parent = match verify_arguments(sess, &mut params, &args) {
        Ok(target) => target,
        Err(e) => {
            println!("ERR {}", e);
            return;
        }
    };

    let vm = rlua::Lua::new();

    let _ = vm.context(|mut vm_context| {
        let _ = vm_context.globals().set("targets", targets.clone());

        let arguments = vm::arg::Args::from(params);
        vm_context.globals().set("args", arguments).unwrap();
        vm_context
            .globals()
            .set("sess_targets", sess.targets.clone())
            .unwrap();

        for extend in &module.metadata.extends {
            vm::extends(&mut vm_context, extend, sess);
        }

        vm::require::common::extends_common(&mut vm_context);

        let content = match std::fs::read(&module.file_name) {
            Ok(c) => c,
            Err(_) => return Err(opf_module::Error::CantLoadContent),
        };

        _ = vm_context.load(content.as_slice()).exec();

        vm_context.scope(|_scope| {
            let new_targets: Vec<HashMap<String, String>> =
                vm_context.globals().get("targets").unwrap_or(vec![]);

            targets = new_targets;
        });
        Ok(())
    });

    let mut inserted = 0;
    for target in targets {
        let mut new_target = match Target::try_from(target) {
            Ok(target) => target,
            Err(e) => {
                println!("ERR target::try_from {}", e);
                continue;
            }
        };

        let module_name = match module.metadata.name.clone() {
            Some(name) => name,
            None => continue,
        };

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

        new_target.target_parent = Some(target_parent.target_uuid);
        new_target.target_groups.push(group.group_uuid);

        if sess.create_target(new_target) {
            inserted += 1;
        }
    }

    println!("! {} new targets", inserted);
}
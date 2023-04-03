mod action;
mod api;
mod export;
mod group;
mod link;
mod module;
mod save;
mod session;
mod target;
pub mod validator;

use crate::core::session::Session;
use crate::error::ErrorKind;
use crate::modules::Manager;
use std::collections::HashMap;
use validator::{CommandAction, CommandObject};

#[derive(Debug)]
pub struct Command {
    pub action: CommandAction,
    object: CommandObject,
    params: HashMap<String, String>,
}

impl Command {
    pub fn get<'a>(&self, name: &str) -> Option<String> {
        if let Some(value) = self.params.get(&name.to_string()) {
            return Some(value.clone());
        }
        None
    }
}

pub fn format(line: &str) -> Result<Command, ErrorKind> {
    let (action, object, params) = {
        let split = line.splitn(3, " ").collect::<Vec<&str>>();
        let split_len = split.len();

        if split_len < 1 {
            return Err(ErrorKind::InvalidCommandArgument);
        }

        let action = validator::validate_action(split[0].to_string())?;

        let mut object = CommandObject::None;
        let mut params = HashMap::new();

        if split_len > 1 {
            object = validator::validate_object(&action, split[1].to_string())?;

            if split_len > 2 {
                let parameters = split[2].split(",").collect::<Vec<&str>>();

                for param in parameters {
                    if !param.contains("=") {
                        return Err(ErrorKind::InvalidFormatArgument);
                    }

                    let vals = param.split("=").collect::<Vec<&str>>();
                    if vals.len() < 1 {
                        return Err(ErrorKind::InvalidFormatArgument);
                    }
                    params.insert(vals[0].trim().to_string(), vals[1].trim().to_string());
                }
            }
        }
        (action, object, params)
    };
    Ok(Command {
        action,
        object,
        params,
    })
}

pub fn exec<'a>(
    session: &'a mut Session,
    cmd: Command,
    module_manager: &'a mut Manager,
) -> Result<(), ErrorKind> {
    match cmd.object {
        CommandObject::None => {}
        CommandObject::Target => target::exec(session, cmd)?,
        CommandObject::Link => link::exec(session, cmd)?,
        CommandObject::Group => group::exec(session, cmd)?,
        CommandObject::Export(_) => export::exec(session, cmd)?,
        CommandObject::Api => {
            println!("API !");
        }
        CommandObject::Module(_) => module::exec(session, cmd, module_manager)?,
        CommandObject::Action => action::exec(session, cmd)?,
        CommandObject::Save => save::exec(session, cmd)?,
        CommandObject::Session => session::exec(session, cmd)?,
    };

    Ok(())
}

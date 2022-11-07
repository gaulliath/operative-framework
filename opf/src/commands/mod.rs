mod error;
mod target;
mod api;
pub mod validator;
mod link;
mod module;
mod export;
mod action;
mod group;
mod session;
mod save;

use crate::modules::Manager;
use crate::core::session::Session;
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

pub fn format(line: &str) -> Result<Command, error::Error> {
    let (action, object, params) = {
        let split = line.splitn(3, " ").collect::<Vec<&str>>();
        let split_len = split.len();

        if split_len < 1 {
            return Err(error::Error::InvalidCommandArgument);
        }

        let action = match validator::validate_action(split[0].to_string()) {
            Ok(action) => action,
            Err(e) => return Err(e),
        };

        let mut object = CommandObject::None;
        let mut params = HashMap::new();

        if split_len > 1 {
            object = match validator::validate_object(&action, split[1].to_string()) {
                Ok(o) => o,
                Err(e) => return Err(e),
            };

            if split_len > 2 {
                let parameters = split[2].split(",").collect::<Vec<&str>>();

                for param in parameters {
                    if !param.contains("=") {
                        return Err(error::Error::InvalidFormatArgument);
                    }

                    let vals = param.split("=").collect::<Vec<&str>>();
                    if vals.len() < 1 {
                        return Err(error::Error::InvalidFormatArgument);
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

pub fn exec<'a>(session: &'a mut Session, cmd: Command, module_manager: &'a mut Manager) -> Result<(), String> {
    match cmd.object {
        CommandObject::None => {},
        CommandObject::Target => {
            if let Err(e) = target::exec(session, cmd) {
                return Err(e.to_string());
            }
        },
        CommandObject::Link => {
            if let Err(e) = link::exec(session, cmd) {
                return Err(e.to_string());
            }
        },
        CommandObject::Group => {
            if let Err(e) = group::exec(session, cmd) {
                return Err(e.to_string())
            }
        }
        CommandObject::Export(_) => {
            if let Err(e) = export::exec(session, cmd) {
                return Err(e.to_string());
            }
        }
        CommandObject::Api => {
            println!("API !");
        }
        CommandObject::Module(_) => {
            if let Err(e) =  module::exec(session, cmd, module_manager) {
                return Err(e.to_string());
            }
        }
        CommandObject::Action => {
            if let Err(e) = action::exec(session, cmd) {
                return Err(e.to_string());
            }
        },
        CommandObject::Save => {
            if let Err(e) = save::exec(session, cmd) {
                return Err(e.to_string());
            }
        }
        CommandObject::Session => {
             if let Err(e) = session::exec(session, cmd) {
                 return Err(e.to_string());
             }
        }
    };

    Ok(())
}

// (opf 2.0) > add target type=company name=seckiot
// > ID = 324342342
// (opf 2.0) > set target id=4324234234 custom_id=seckiot_1
// (opf 2.0) > run linkedin.employee target=seckiot_1
// (opf : linkedin.employee) > set target seckiot_1
// (opf : linkedin.employee) > set output false && run

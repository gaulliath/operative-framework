pub mod link;
pub mod output;
pub mod session;
pub mod target;
pub mod group;

use crate::commands;
use crate::modules::Manager;
use crate::modules::Module;
use session::{Session, SessionConfig};
use crate::error::Manager as ManagerError;

#[derive(Debug)]
pub struct Core {
    pub session: Session,
    config: SessionConfig,
    commands: Vec<String>,
    manager: Manager,
}

pub fn new(session: Session, config: SessionConfig) -> Core {
    Core {
        session,
        config,
        commands: vec![],
        manager: Manager::default(),
    }
}

impl Core {
    pub fn session_id(&self) -> String {
        format!("{}", &self.session.session_id)
    }

    pub fn get_session(&self) -> &Session {
        &&self.session
    }

    pub fn init_manager<'a>(&mut self, path: &'a str) -> Result<(), ManagerError> {
        crate::modules::register_modules(&mut self.manager, path)
    }

    pub fn get_modules(&self) -> Vec<Module> {
        return self.manager.modules.clone()
    }

    pub fn run(&mut self, ask: String) -> Result<(), String> {
        let command = match commands::format(ask.clone().as_str()) {
            Ok(command) => command,
            Err(e) => {
                return Err(e.to_string());
            }
        };

        if self.config.verbose {
            println!("{:?}", command);
        }

        let module_manager = &mut self.manager;
        if command.action.eq(&commands::validator::CommandAction::Run) {}

        if let Err(e) = commands::exec(&mut self.session, command, module_manager) {
            return Err(e);
        }

        self.commands.push(ask);
        Ok(())
    }

    pub fn commands(&self) -> Vec<String> {
        self.commands.clone()
    }
}

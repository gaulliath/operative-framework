use session::{Session, SessionConfig};

use crate::commands;
use crate::error::ErrorKind;
use crate::modules::Manager;
use crate::modules::Module;

pub mod group;
pub mod link;
pub mod output;
pub mod session;
pub mod target;

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

    pub fn init_manager(&mut self, path: &str) -> Result<(), ErrorKind> {
        crate::modules::register_modules(&mut self.manager, path)
    }

    pub fn get_modules(&self) -> Vec<Module> {
        return self.manager.modules.clone();
    }

    pub fn run(&mut self, ask: &str) -> Result<(), ErrorKind> {
        let command = commands::format(ask)?;

        if self.config.verbose {
            println!("{:?}", command);
        }

        let module_manager = &mut self.manager;
        if command.action.eq(&commands::validator::CommandAction::Run) {}

        commands::exec(&mut self.session, command, module_manager)?;

        self.commands.push(ask.to_string());
        Ok(())
    }

    pub fn commands(&self) -> Vec<String> {
        self.commands.clone()
    }
}

use session::{Session, SessionConfig};

use crate::commands;
use crate::error::ErrorKind;
use crate::modules::Manager;

pub mod group;
pub mod link;
pub mod output;
pub mod session;
pub mod target;

#[derive(Debug)]
#[allow(dead_code)]
pub struct Core {
    pub session: Session,
    config: SessionConfig,
    manager: Manager,
}

pub fn new(session: Session, config: SessionConfig) -> Core {
    Core {
        session,
        config,
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
        crate::compiled::registers(
            &mut self.manager,
            vec![
                crate::compiled::sample::SampleCompiled::new(),
                crate::compiled::linkedin_search::LinkedinSearch::new(),
            ],
        )?;
        crate::modules::register(&mut self.manager, path)?;
        Ok(())
    }

    pub fn run(&mut self, ask: &str) -> Result<(), ErrorKind> {
        commands::exec(&mut self.session, commands::format(ask)?, &mut self.manager)?;
        Ok(())
    }
}

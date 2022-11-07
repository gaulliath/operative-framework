use super::{Command, CommandAction, CommandObject};
use crate::core::session::Session;
use crate::common::module as opf_module;

pub fn exec<'a>(
    session: &'a mut Session,
    cmd: Command,
) -> Result<(), opf_module::Error> {
    match cmd.action {
        CommandAction::Run => run(session, cmd),
        _ => Ok(()),
    }
}


pub fn run<'a>(
    session: &'a mut Session,
    cmd: Command,
) -> Result<(), opf_module::Error> {

    Ok(())
}


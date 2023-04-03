#![allow(dead_code)]
use super::{Command, CommandAction};
use crate::common::module as opf_module;
use crate::core::session::Session;

pub fn exec(session: &mut Session, cmd: Command) -> Result<(), opf_module::Error> {
    match cmd.action {
        CommandAction::Run => run(session, cmd),
        _ => Ok(()),
    }
}

pub fn run(_session: &mut Session, _cmd: Command) -> Result<(), opf_module::Error> {
    Ok(())
}

use super::{Command, CommandAction};
use crate::core::session::Session;
use crate::error::ErrorKind;
use colored::*;
use crossterm::style::Stylize;

pub fn exec<'a>(_session: &'a mut Session, cmd: Command) -> Result<(), ErrorKind> {
    match cmd.action {
        CommandAction::Help => help(),
        _ => Ok(()),
    }
}

fn help() -> Result<(), ErrorKind> {
    println!(
        "{} : save session in specified file",
        "save session to=<path>".bright_yellow()
    );
    println!(
        "{}",
        "- save session to=/tmp/session.today : export session into /tmp/session.today file".grey()
    );
    Ok(())
}

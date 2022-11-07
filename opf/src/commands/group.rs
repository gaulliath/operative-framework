use crate::commands::{Command, CommandAction};
use crate::core::session::Session;
use crate::common::groups as opf_group;

use colored::*;
use crossterm::style::Stylize;

pub fn exec<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_group::Error> {
    match cmd.action {
        CommandAction::List => list(session, cmd),
        CommandAction::Help => help(),
        _ => Ok(()),
    }
}

fn help() -> Result<(), opf_group::Error> {
    println!(
        "{} : show all group available in current session",
        "list group".bright_yellow()
    );
    println!(
        "{}",
        "- list group metadata=true: this command print all group available with metadata."
            .grey()
    );
    Ok(())
}


fn list(session: &mut Session, cmd: Command) -> Result<(), opf_group::Error> {
    let headers = vec!["uuid".to_string(), "uid".to_string(), "name".to_string()];
    let mut rows = vec![];

    for group in &session.groups {
        rows.push(vec![
            group.group_uuid.clone().to_string(),
            group.group_id.clone().to_string(),
            group.group_name.clone(),
        ]);
    }

    session.output_table(headers, rows);
    Ok(())
}

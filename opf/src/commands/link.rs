use super::{Command, CommandAction};
use crate::common::search::Link as link_search;
use crate::core::session::Session;
use crate::error::{ErrorKind, Link as LinkError};
use colored::*;
use crossterm::style::Stylize;
use std::str::FromStr;
use std::time::SystemTime;

pub fn exec(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    match cmd.action {
        CommandAction::Add => add(session, cmd),
        CommandAction::Del => del(session, cmd),
        CommandAction::Set => set(session, cmd),
        CommandAction::List => list(session, cmd),
        CommandAction::Help => help(),
        _ => Ok(()),
    }
}

fn help() -> Result<(), ErrorKind> {
    println!(
        "{}: this command add a new link into session.",
        "add link".bright_yellow()
    );
    println!(
        "{}",
        "- example: add link src=target_id, dst=target_id, metadata_1=text, metadata_2=text".grey()
    );
    println!(
        "{}: remove a link from current session.",
        "del link".bright_yellow()
    );
    println!("{}", "- example: del link id=id_of_link");
    println!(
        "{}: update information or metadata for specified link.",
        "set link".bright_yellow()
    );
    println!(
        "{}",
        "- set link id=id_of_link, my_new_metadata=value".grey()
    );
    println!(
        "{} : show all link available in current session",
        "list link".bright_yellow()
    );
    println!(
        "{}",
        "- list link metadata=true: this command print all link available with metadata.".grey()
    );
    Ok(())
}

fn list(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;
    let search = link_search::from(params);

    let links = &session.get_links(search);

    let headers = vec![
        "id".to_string(),
        "source".to_string(),
        "target".to_string(),
        "label".to_string(),
        "color".to_string(),
    ];
    let mut rows = vec![];

    for link in links {
        rows.push(vec![
            link.link_id.to_string(),
            link.link_target.to_string(),
            link.link_source.to_string(),
            link.link_label.to_string(),
            link.link_color.clone().unwrap_or(String::from("-")),
        ]);
    }

    session.output_table(headers, rows);
    Ok(())
}

fn add(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;

    let link_label = match params.get("label") {
        Some(label) => label.clone(),
        None => {
            return Err(ErrorKind::Link(LinkError::ParamNotFound(
                "label".to_string(),
            )))
        }
    };

    let link_color = match params.get("color") {
        Some(color) => Some(color.clone()),
        None => None,
    };

    let link_created_by = match params.get("created_by") {
        Some(created_by) => created_by.clone(),
        None => {
            return Err(ErrorKind::Link(LinkError::ParamNotFound(
                "created_by".to_string(),
            )))
        }
    };

    let link_source = match params.get("source") {
        Some(source) => match uuid::Uuid::from_str(source) {
            Ok(uuid) => uuid,
            Err(_) => {
                return Err(ErrorKind::Link(LinkError::ParamFormatInvalid(
                    "source".to_string(),
                )));
            }
        },
        None => {
            return Err(ErrorKind::Link(LinkError::ParamNotFound(
                "source".to_string(),
            )))
        }
    };

    let link_target = match params.get("target") {
        Some(target) => match uuid::Uuid::from_str(target) {
            Ok(uuid) => uuid,
            Err(_) => {
                return Err(ErrorKind::Link(LinkError::ParamFormatInvalid(
                    "target".to_string(),
                )));
            }
        },
        None => {
            return Err(ErrorKind::Link(LinkError::ParamNotFound(
                "target".to_string(),
            )))
        }
    };

    let search_link = link_search {
        link_id: None,
        link_label: Some(link_label.clone()),
        link_color: None,
        link_source: Some(link_source.clone().to_string()),
        link_target: Some(link_target.clone().to_string()),
        link_created_by: Some(link_created_by.clone()),
    };

    if session.exist_link(search_link) {
        return Err(ErrorKind::Link(LinkError::LinkExist));
    }

    let new_link = crate::common::Link {
        link_id: uuid::Uuid::new_v4(),
        link_label,
        link_color,
        link_source,
        link_target,
        link_created_by,
        link_created_at: SystemTime::now(),
    };

    println!("{:#?}", new_link);
    session.create_link(new_link);
    Ok(())
}

fn del(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;
    let search = link_search::from(params);
    session.delete_links(search);
    println!("links available -> ({})", session.links.len());
    Ok(())
}

fn set(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;
    let search = link_search::from(params);
    _ = session.update_links(search);
    Ok(())
}

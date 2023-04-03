use std::str::FromStr;

use super::{Command, CommandAction};
use crate::common::search::Target as target_search;
use crate::common::target as opf_target;
use crate::core::session::Session;
use crate::error::{ErrorKind, Target as TargetError};
use crate::utils;
use colored::*;
use crossterm::style::Stylize;

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
        "{}: this command add a new target into session.",
        "add target".bright_yellow()
    );
    println!(
        "{}",
        "- example: add target type=company, name=value".grey()
    );
    println!(
        "{}: remove a target from current session.",
        "del target".bright_yellow()
    );
    println!("{}", "- example: del target id=id_of_target");
    println!(
        "{}: update information or metadata for specified target.",
        "set target".bright_yellow()
    );
    println!(
        "{}",
        "- set target id=id_of_target, my_new_metadata=value".grey()
    );
    println!(
        "{} : show all target available in current session",
        "list target".bright_yellow()
    );
    println!(
        "{}",
        "- list target metadata=true: this command print all target available with metadata."
            .grey()
    );
    Ok(())
}

fn list(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;
    let search = target_search::from(params.clone());
    let targets = &session.get_targets(search);

    let show_metadata = params.get("metadata").is_some();

    let mut headers = vec![
        "uuid".to_string(),
        "uid".to_string(),
        "parent".to_string(),
        "custom_id".to_string(),
        "type".to_string(),
        "name".to_string(),
    ];
    let mut header_meta = vec![];
    let mut rows = vec![];

    for target in targets {
        let parent = match target.target_parent {
            Some(id) => id.to_string(),
            None => String::from("-"),
        };

        let mut fields = vec![
            target.target_uuid.to_string(),
            target.target_id.to_string(),
            parent,
            target.target_custom_id.clone().unwrap_or(String::from("-")),
            target.target_type.to_string(),
            target.target_name.clone(),
        ];

        if show_metadata {
            for (key, _) in &target.meta.clone() {
                let key_s = String::from(key);
                if header_meta.contains(&key_s) {
                    continue;
                }
                header_meta.push(key_s);
            }

            for head in &header_meta {
                match target.meta.get(head) {
                    Some(meta) => {
                        fields.push(meta.clone());
                    }
                    None => fields.push("-".to_string()),
                }
            }
        }
        rows.push(fields);
    }

    headers.append(&mut header_meta);
    session.output_table(headers, rows);
    Ok(())
}

fn add(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let mut params = cmd.params;

    let target_name = match params.get("name") {
        Some(name) => name.clone(),
        None => return Err(ErrorKind::Target(TargetError::ParamNameNotFound)),
    };

    let target_type = match params.get("type") {
        Some(t) => opf_target::validate_type(t)?,
        None => return Err(ErrorKind::Target(TargetError::ParamTypeNotFound)),
    };

    params.remove("name");
    params.remove("type");

    let target_parent = match params.get("parent") {
        Some(parent_uuid) => {
            let mut search = target_search::default();
            search.target_uuid = Some(parent_uuid.clone());
            if !session.exist_target(&search) {
                return Err(ErrorKind::Target(TargetError::ParentUuidNotFound));
            }

            match uuid::Uuid::from_str(&parent_uuid) {
                Ok(uuid) => Some(uuid),
                Err(_) => {
                    return Err(ErrorKind::Target(TargetError::ParentUuidNotValid));
                }
            }
        }
        None => None,
    };

    let target_custom_id = match utils::param::get_param(&params, "custom_id") {
        Ok(p) => Some(p),
        Err(_) => None,
    };

    let mut target_groups = vec![];
    if let Ok(groups) = utils::param::get_param(&params, "groups") {
        if groups.contains(",") {
            let groups: Vec<&str> = groups.split(",").collect();
            for group in groups {
                if let Ok(group) = session.get_group(group) {
                    target_groups.push(group.group_uuid);
                }
            }
        } else {
            if let Ok(group) = session.get_group(groups.as_str()) {
                target_groups.push(group.group_uuid);
            }
        }
    }

    let target_parent_string = match target_parent.clone() {
        Some(uuid) => Some(uuid.to_string()),
        None => None,
    };

    let search_target = target_search {
        target_uuid: None,
        target_id: None,
        target_parent: target_parent_string,
        target_custom_id: None,
        target_group_uuid: None,
        target_name: Some(target_name.clone()),
        target_type: Some(target_type.clone()),
        meta: None,
    };

    if session.exist_target(&search_target) {
        return Err(ErrorKind::Target(TargetError::TargetExist));
    }

    let target_id = session.targets.len() as i32 + 1;

    let new_target = crate::common::Target {
        target_uuid: uuid::Uuid::new_v4(),
        target_id,
        target_name,
        target_type,
        target_custom_id,
        target_groups,
        target_parent,
        meta: params.clone(),
    };

    println!(
        "! id={}, uuid={}",
        new_target.target_id, new_target.target_uuid
    );
    let _ = session.create_target(new_target)?;
    Ok(())
}

fn del(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let params = cmd.params;
    let search = target_search::from(params);
    session.delete_targets(search);
    println!("targets available -> ({})", session.targets.len());
    Ok(())
}

fn set(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    println!("{:#?}", cmd);
    let params = cmd.params;
    let search = target_search::from(params);
    _ = session.update_targets(search);
    Ok(())
}

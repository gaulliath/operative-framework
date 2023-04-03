use super::{Command, CommandAction, CommandObject};
use crate::common::export as opf_export;
use crate::core::session::Session;
use crate::error::{ErrorKind, Export as ExportError};
use colored::*;
use crossterm::style::Stylize;
use petgraph::dot::Dot;
use petgraph::graph::Graph;
use std::collections::HashMap;
use std::str::FromStr;

pub fn exec(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    match cmd.action {
        CommandAction::Export => export(session, cmd),
        CommandAction::Help => help(),
        _ => Ok(()),
    }
}

fn help() -> Result<(), ErrorKind> {
    println!(
        "{} : export session in specified type",
        "export <type>".bright_yellow()
    );
    println!(
        "{}",
        "- export dot : this command export session in dot graph format.".grey()
    );
    Ok(())
}

pub fn export(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let export_type = match cmd.object {
        CommandObject::Export(export_type) => export_type,
        _ => return Err(ErrorKind::Export(ExportError::ExportType)),
    };

    let export_type = match opf_export::ExportType::from_str(export_type.as_str()) {
        Ok(t) => t,
        Err(_) => return Err(ErrorKind::Export(ExportError::ExportType)),
    };

    match export_type {
        opf_export::ExportType::Dot => export_dot(session, &cmd.params),
    }

    Ok(())
}

fn export_dot(session: &Session, _params: &HashMap<String, String>) {
    let mut graph = Graph::<String, String>::new();
    let mut index = HashMap::new();

    for target in &session.targets {
        let id = graph.add_node(target.target_name.clone());
        index.insert(target.target_uuid.clone(), id.clone());

        if let Some(parent) = target.target_parent {
            if let Some(parent_id) = index.get(&parent) {
                graph.add_edge(
                    parent_id.clone(),
                    id,
                    target.target_type.to_string().clone(),
                );
            }
        }
    }

    println!("{:?}", Dot::new(&graph));
}

use super::{Command, CommandAction, CommandObject};
use crate::core::session::Session;
use crate::common::export as opf_export;
use std::str::FromStr;
use colored::*;
use crossterm::style::Stylize;
use std::collections::HashMap;
use petgraph::dot::Dot;
use petgraph::graph::Graph;

pub fn exec<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_export::Error> {
    match cmd.action {
        CommandAction::Export => export(session, cmd),
        CommandAction::Help => help(),
        _ => Ok(()),
    }
}

fn help() -> Result<(), opf_export::Error> {
    println!(
        "{} : export session in specified type",
        "export <type>".bright_yellow()
    );
    println!(
        "{}",
        "- export dot : this command export session in dot graph format."
            .grey()
    );
    Ok(())
}

pub fn export<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_export::Error> {
    let export_type = match cmd.object {
        CommandObject::Export(export_type) => export_type,
        _ => return Err(opf_export::Error::ExportType),
    };

    let export_type = match opf_export::ExportType::from_str(export_type.as_str()) {
        Ok(t) => t,
        Err(_) => return Err(opf_export::Error::ExportType),
    };

    match export_type {
        opf_export::ExportType::Dot => export_dot(session, &cmd.params),
    }

    Ok(())
}

fn export_dot<'a, 'b>(session: &'a Session, _params: &'b HashMap<String, String>) {
    let mut graph = Graph::<String, String>::new();
    let mut index = HashMap::new();

    for target in &session.targets {
        let id = graph.add_node(target.target_name.clone());
        index.insert(target.target_uuid.clone(), id.clone());

        if let Some(parent) = target.target_parent {
            if let Some(parent_id) = index.get(&parent) {
                graph.add_edge(parent_id.clone(), id, target.target_type.to_string().clone());
            }
        }
    }

    println!("{:?}", Dot::new(&graph));
}

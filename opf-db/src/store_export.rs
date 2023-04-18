use opf_models::event::{send_event, Event};
use opf_models::{self, error::ErrorKind, Command, CommandAction};
use petgraph::dot::Dot;
use petgraph::graph::Graph;
use std::collections::HashMap;

use crate::store::DB;

impl DB {
    pub async fn on_export_command(&mut self, command: Command) -> Result<(), ErrorKind> {
        match command.action {
            CommandAction::Export => self.export_dot(command).await,
            _ => Err(ErrorKind::ActionNotAvailable),
        }
    }

    async fn export_dot(&mut self, _cmd: Command) -> Result<(), ErrorKind> {
        let mut graph = Graph::<String, String>::new();
        let mut index = HashMap::new();
        let targets = self.targets.read().await;
        for (target_id, target) in targets.iter() {
            if target.target_parent.is_some() {
                continue;
            }
            let id = graph.add_node(target.target_name.clone());
            index.insert(target_id, id.clone());
        }

        for (target_id, target) in targets.iter() {
            if target.target_parent.is_none() {
                continue;
            }
            let id = graph.add_node(target.target_name.clone());
            index.insert(target_id, id.clone());

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

        send_event(
            &self.db_tx,
            Event::ResponseSimple(format!("{:?}", Dot::new(&graph))),
        )
        .await
    }
}

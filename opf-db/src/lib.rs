use std::collections::HashMap;
use std::sync::Arc;
use std::time::SystemTime;

use log::info;
use rand::{distributions::Alphanumeric, Rng};
use tokio::sync::mpsc::{unbounded_channel, UnboundedReceiver, UnboundedSender};
use tokio::sync::RwLock;

use opf_models::event::Event::{ResponseError, ResponseSimple};
use opf_models::event::{send_event_to, Domain, Event};
use opf_models::{KeyStore, Workspace};

use crate::store::DB;

mod store;
mod store_export;
mod store_group;
mod store_keystore;
mod store_link;
mod store_module;
mod store_target;
mod store_workspace;

#[derive(Debug)]
pub struct DBStore {
    dbs: HashMap<i32, DB>,
    keystore: Arc<RwLock<KeyStore>>,
    workspaces: HashMap<i32, Workspace>,
    pub current_workspace: i32,
    pub self_tx: UnboundedSender<Event>,
    pub self_rx: UnboundedReceiver<Event>,
    pub node_tx: UnboundedSender<(Domain, Event)>,
}

pub async fn new(node_tx: UnboundedSender<(Domain, Event)>) -> (UnboundedSender<Event>, DBStore) {
    let (self_tx, self_rx) = unbounded_channel::<Event>();

    // generate a first workspace name, can be update.
    let workspace_id = 0;
    let random_name: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(7)
        .map(char::from)
        .collect();
    let workspace_name = format!("tmp/{}", random_name);

    let temp_db = DB::new(self_tx.clone());
    let mut dbs = HashMap::new();
    dbs.insert(workspace_id, temp_db);
    let tx = self_tx.clone();

    let _ = send_event_to(
        &node_tx,
        (Domain::CLI, Event::SetWorkspace(workspace_name.clone())),
    )
    .await;

    let mut workspaces = HashMap::new();
    workspaces.insert(
        workspace_id,
        Workspace {
            workspace_id,
            workspace_name,
            workspace_created_at: SystemTime::now(),
        },
    );

    (
        tx,
        DBStore {
            keystore: Arc::new(RwLock::new(KeyStore(HashMap::new()))),
            dbs,
            workspaces,
            current_workspace: 0,
            self_tx,
            self_rx,
            node_tx,
        },
    )
}

impl DBStore {
    pub async fn launch(mut self) {
        info!("running database controller...");
        loop {
            tokio::select! {
                Some(event) = self.self_rx.recv() => {
                    let db = match self.dbs.get_mut(&self.current_workspace) {
                        Some(db) => db,
                        None => {
                            let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseSimple(format!("workspace id = {} not found...", self.current_workspace)))).await;
                            continue;
                        }
                    };

                    match event {
                        Event::CommandTarget(command) => {
                            if let Err(e) = db.on_target_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        },
                        Event::CommandExport(command) => {
                            if let Err(e) = db.on_export_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        },
                        Event::CommandGroup(command) => {
                            if let Err(e) = db.on_group_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        },
                        Event::PrepareModule((module_name, command)) => {
                            let keystore = self.keystore.read().await;
                            if let Err(e) = db.on_prepare_module(&keystore, module_name, command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::UpdateTargetMeta((target_id, meta)) => {
                            if let Err(e) = db.on_target_meta_update(target_id, meta).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::ResultsModule(group_id, targets) => {
                            if let Err(e) = db.on_results_targets(group_id, targets).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::ExecuteModule(_) => {
                             let _ = send_event_to(&self.node_tx, (Domain::Module, event)).await;
                        }
                        Event::CommandLink(command) => {
                            if let Err(e) = db.on_link_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::CommandWorkspace(command) => {
                            if let Err(e) = self.on_workspace_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::CommandKeystore(command) => {
                            if let Err(e) = self.on_keystore_command(command).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::LoadKeystore(keystore) => self.load_keystore(keystore).await,
                        _ => {
                            if let Err(e) = send_event_to(&self.node_tx, (Domain::CLI, event)).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                    }
                }
            }
        }
    }

    pub async fn load_keystore(&mut self, keystore: KeyStore) {
        self.keystore = Arc::new(RwLock::new(keystore));
    }
}

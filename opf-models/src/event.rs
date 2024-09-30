use std::collections::HashMap;

use tokio::sync::mpsc::UnboundedSender;

use crate::error::ErrorKind;
use crate::{Command, KeyStore, Target};

#[derive(Debug)]
pub enum Event {
    // command variants
    NewCommand(String),
    ProcessCommand(Command),
    CommandTarget(Command),
    CommandLink(Command),
    CommandModule(Command),
    CommandWorkspace(Command),
    CommandKeystore(Command),
    CommandExport(Command),
    CommandGroup(Command),
    // Module variant
    PrepareModule((String, Command)),
    ExecuteModule((i32, Target, String, HashMap<String, String>)),
    ListModules,
    HelpModule(String),
    ResultsModule(i32, Vec<Target>),
    // response variants
    ResponseSimple(String),
    ResponseError(String),
    ResponseInfo(String),
    ResponseTable((Vec<String>, Vec<Vec<String>>)),
    SetWorkspace(String),
    // keystore
    LoadKeystore(KeyStore),
    // target
    UpdateTargetMeta((String, (String, String))),
}

#[derive(Debug)]
pub enum Domain {
    Data,
    Module,
    Node,
    Network,
    CLI,
}

pub async fn send_event(tx: &UnboundedSender<Event>, event: Event) -> Result<(), ErrorKind> {
    tx.send(event)
        .map_err(|e| ErrorKind::Channel(e.to_string()))
}

pub async fn send_event_to(
    tx: &UnboundedSender<(Domain, Event)>,
    message: (Domain, Event),
) -> Result<(), ErrorKind> {
    tx.send(message)
        .map_err(|e| ErrorKind::Channel(e.to_string()))
}

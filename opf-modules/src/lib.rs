use std::collections::HashMap;

use async_trait::async_trait;
use dyn_clone::DynClone;
use log::info;
use tokio::sync::mpsc::{unbounded_channel, UnboundedReceiver, UnboundedSender};

use opf_models::event::send_event_to;
use opf_models::event::Event::{ResponseError, ResultsModule};
use opf_models::metadata::{Arg, Args, Metadata};
use opf_models::{
    error::ErrorKind,
    event::{Domain, Event},
    ModuleAction, Target, TargetType,
};

mod account_search;
mod command;
mod crtsh;
mod email_to_domain;
mod insee_search;
mod linkedin_search;
mod port_scanner;
mod worker;

pub enum Module {
    Lua(LuaModule),
    Compiled(Box<dyn CompiledModule>),
}

#[derive(Debug, Clone)]
pub struct LuaModule {
    pub file_name: String,
    pub metadata: Metadata,
}

pub struct ModuleController {
    pub self_tx: UnboundedSender<Event>,
    pub self_rx: UnboundedReceiver<Event>,
    pub node_tx: UnboundedSender<(Domain, Event)>,
    pub modules: HashMap<String, Module>,
}

#[async_trait]
pub trait CompiledModule: Sync + Send + DynClone {
    fn name(&self) -> String;
    fn author(&self) -> String;
    fn resume(&self) -> String;
    fn args(&self) -> Vec<Arg>;
    fn target_type(&self) -> TargetType;
    fn module_action(&self) -> ModuleAction {
        ModuleAction::CreateTarget
    }
    fn is_threaded(&self) -> bool {
        false
    }
    async fn run(
        &self,
        params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind>;
}

impl Module {
    pub fn name(&self) -> String {
        match &self {
            Self::Lua(ref module) => module.metadata.name.clone().unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.name(),
        }
    }

    pub fn author(&self) -> String {
        match &self {
            Self::Lua(ref module) => module.metadata.author.clone().unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.author(),
        }
    }

    pub fn resume(&self) -> String {
        match &self {
            Self::Lua(ref module) => module
                .metadata
                .description
                .clone()
                .unwrap_or("-".to_string()),
            Self::Compiled(ref module) => module.resume(),
        }
    }

    pub fn args(&self) -> Vec<Arg> {
        match &self {
            Self::Lua(ref module) => module.metadata.args.clone(),
            Self::Compiled(ref module) => module.args(),
        }
    }

    pub fn target_type(&self) -> TargetType {
        match &self {
            Self::Lua(ref _module) => TargetType::Person,
            Self::Compiled(ref module) => module.target_type(),
        }
    }
}

impl ModuleController {
    pub async fn launch(mut self) {
        info!("running module controller...");
        loop {
            tokio::select! {
                Some(event) = self.self_rx.recv() => {
                    match event {
                        Event::ExecuteModule(data) => {
                            if let Err(e) = self.on_execute_module(data).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::ListModules => {
                            if let Err(e) = self.on_list_modules().await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::HelpModule(module_name) => {
                            if let Err(e) = self.on_help_module(module_name).await {
                                let _ = send_event_to(&self.node_tx, (Domain::CLI, ResponseError(e.to_string()))).await;
                            }
                        }
                        Event::ResultsModule(targets) => {
                            let _ = send_event_to(&self.node_tx, (Domain::Data, ResultsModule(targets))).await;
                        }
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
}

pub fn new(
    node_tx: UnboundedSender<(Domain, Event)>,
) -> (UnboundedSender<Event>, ModuleController) {
    let (self_tx, self_rx) = unbounded_channel::<Event>();
    let tx = self_tx.clone();

    let mut modules = HashMap::new();
    let linkedin = linkedin_search::LinkedinSearch::new();
    let account_search = account_search::AccountSearch::new();
    let port_scanner = port_scanner::PortScanner::new();
    let crt_sh = crtsh::CrtSH::new();
    let email_to_domain = email_to_domain::EmailToDomain::new();
    let insee_search = insee_search::InseeSearch::new();

    modules.insert(linkedin.name(), Module::Compiled(linkedin));
    modules.insert(account_search.name(), Module::Compiled(account_search));
    modules.insert(port_scanner.name(), Module::Compiled(port_scanner));
    modules.insert(crt_sh.name(), Module::Compiled(crt_sh));
    modules.insert(email_to_domain.name(), Module::Compiled(email_to_domain));
    modules.insert(insee_search.name(), Module::Compiled(insee_search));

    (
        tx,
        ModuleController {
            self_tx,
            self_rx,
            node_tx,
            modules,
        },
    )
}

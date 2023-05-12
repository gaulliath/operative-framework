use async_trait::async_trait;
use fantoccini::Locator;
use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use tokio::sync::mpsc::UnboundedSender;

use crate::CompiledModule;

#[derive(Clone)]
pub struct EnterpriseDomain {}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Res {
    pub suggestions: Vec<String>,
}

#[async_trait]
impl CompiledModule for EnterpriseDomain {
    fn name(&self) -> String {
        "enterprise.domain".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Search for a domain name associated with an enterprise.".to_string()
    }

    fn is_threaded(&self) -> bool {
        true
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    fn target_type(&self) -> TargetType {
        TargetType::Company
    }

    async fn run(
        &self,
        group_id: i32,
        target: Target,
        _params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let target_id = target.target_id.to_string();
        let tx = tx.unwrap();
        let url = format!(
            "https://securitytrails.com/api/public/app/api/vercel/autocomplete/domain/{}",
            target.target_name
        );
        let client = fantoccini::ClientBuilder::native()
            .connect("http://localhost:4444")
            .await
            .map_err(|e| ErrorKind::GenericError(format!("web driver : {}", e)))?;

        client
            .goto(&url)
            .await
            .map_err(|e| ErrorKind::GenericError(format!("web driver : {}", e)))?;

        let src = client
            .find(Locator::Css("#json"))
            .await
            .map_err(|e| ErrorKind::GenericError(format!("source error : {}", e)))?;

        let source = src
            .html(true)
            .await
            .map_err(|e| ErrorKind::GenericError(format!("source error : {}", e)))?;

        let parsing = serde_json::from_str::<Res>(source.as_str()).map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(format!("parsing of json error.")))
        })?;

        let mut results = vec![];
        for domain in parsing.suggestions {
            let mut target = HashMap::new();
            target.insert(String::from(opf_models::target::NAME), domain);
            target.insert(
                String::from(opf_models::target::TYPE),
                TargetType::Domain.to_string(),
            );
            target.insert(String::from(opf_models::target::PARENT), target_id.clone());
            results.push(Target::try_from(target)?);
        }

        let _ = send_event(&tx, Event::ResultsModule(group_id, results)).await;
        Ok(vec![])
    }
}

impl EnterpriseDomain {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(EnterpriseDomain {})
    }
}

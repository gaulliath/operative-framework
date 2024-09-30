use async_trait::async_trait;
use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};
use reqwest::header;
use tokio::sync::mpsc::UnboundedSender;

use crate::CompiledModule;

#[derive(Clone)]
pub struct DomainAlive {}

#[async_trait]
impl CompiledModule for DomainAlive {
    fn name(&self) -> String {
        "domain.alive".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Check if domain is alive".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    fn target_type(&self) -> TargetType {
        TargetType::Domain
    }

    async fn run(
        &self,
        _group_id: i32,
        target: Target,
        _params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let target_id = target.target_id.to_string();

        let https_url = String::from(format!("https://{}", target.target_name.clone()));
        let http_url = String::from(format!("http://{}", target.target_name.clone()));

        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            header::HeaderValue::from_static(
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
            ),
        );

        let client = reqwest::Client::builder()
            .default_headers(headers)
            .build()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let tx = tx.unwrap();
        match client.get(https_url).send().await {
            Ok(_) => {
                let _ = send_event(
                    &tx,
                    Event::UpdateTargetMeta((
                        target_id.clone(),
                        (opf_models::domain::IS_ALIVE.to_string(), "true".to_string()),
                    )),
                )
                .await;

                return Ok(vec![]);
            }
            Err(_) => {
                let _ = send_event(
                    &tx,
                    Event::UpdateTargetMeta((
                        target_id.clone(),
                        (
                            opf_models::domain::IS_ALIVE.to_string(),
                            "false".to_string(),
                        ),
                    )),
                )
                .await;
            }
        }

        match client.get(http_url).send().await {
            Ok(_) => {
                let _ = send_event(
                    &tx,
                    Event::UpdateTargetMeta((
                        target_id.clone(),
                        (opf_models::domain::IS_ALIVE.to_string(), "true".to_string()),
                    )),
                )
                .await;
            }
            Err(_) => {
                let _ = send_event(
                    &tx,
                    Event::UpdateTargetMeta((
                        target_id.clone(),
                        (
                            opf_models::domain::IS_ALIVE.to_string(),
                            "false".to_string(),
                        ),
                    )),
                )
                .await;
            }
        }

        Ok(vec![])
    }
}

impl DomainAlive {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(DomainAlive {})
    }
}

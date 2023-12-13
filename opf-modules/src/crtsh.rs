use async_trait::async_trait;
use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::Event;
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};
use reqwest::header;
use scraper::{Html, Selector};
use std::collections::HashMap;
use std::ops::Add;
use tokio::sync::mpsc::UnboundedSender;

use crate::CompiledModule;

#[derive(Clone)]
pub struct CrtSH {}

#[async_trait]
impl CompiledModule for CrtSH {
    fn name(&self) -> String {
        "crt.sh".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "search (sub)domain name by SSL certificate".to_string()
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
        _tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let target_id = target.target_id.to_string();
        let target = target.target_name.clone();

        let url = String::from("https://crt.sh/?q=").add(target.replace(" ", "+").as_str());

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

        let res = client
            .get(url)
            .send()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let res = res
            .text()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let fragment = Html::parse_fragment(res.as_str());
        let selector = Selector::parse(".outer tr>td").map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(
                "invalid selector from scrapping".to_string(),
            ))
        })?;
        let mut results = vec![];
        let mut tmp_results = vec![];

        for element in fragment.select(&selector) {
            let html = element.inner_html();
            if html.contains(&target) {
                if html.contains(&"<br>") {
                    for res in html.split(&"<br>") {
                        let res = res.trim().to_string();
                        if !tmp_results.contains(&res) {
                            tmp_results.push(res);
                        }
                    }
                    continue;
                }

                let res = html.trim().to_string();
                if !tmp_results.contains(&res) {
                    tmp_results.push(res);
                }
            }
        }

        for domain in tmp_results {
            let mut target = HashMap::new();
            target.insert(String::from(opf_models::target::NAME), domain);
            target.insert(
                String::from(opf_models::target::TYPE),
                TargetType::Domain.to_string(),
            );
            target.insert(String::from(opf_models::target::PARENT), target_id.clone());
            results.push(Target::try_from(target)?);
        }
        Ok(results)
    }
}

impl CrtSH {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(CrtSH {})
    }
}

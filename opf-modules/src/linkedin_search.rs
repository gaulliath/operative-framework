use std::collections::HashMap;
use std::ops::Add;

use async_trait::async_trait;

use reqwest::header;
use scraper::{Html, Selector};
use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::Event;
use opf_models::{
    metadata::{Arg, Args},
    Target,
};

use crate::CompiledModule;

#[derive(Clone)]
pub struct LinkedinSearch {}

#[async_trait]
impl CompiledModule for LinkedinSearch {
    fn name(&self) -> String {
        "linkedin.search".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "Search employee on selected enterprise with Linkedin".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
            Arg::new("limit", false, true, Some("10".to_string())),
        ]
    }

    async fn run(
        &self,
        params: Args,
        _tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let company = params.get("target").unwrap();
        let target_id = params.get("target_id").unwrap();
        let limit = params.get("limit").unwrap();
        let target = company.value.unwrap();

        let url = String::from("https://www.google.com/search?num=")
            .add(&limit.value.unwrap().as_str())
            .add("&start=0&hl=en&q=site:linkedin.com/in+")
            .add(&target.replace(" ", "+").as_str());

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
        let selector = Selector::parse(".g h3").map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(
                "invalid selector from scrapping".to_string(),
            ))
        })?;
        let mut results = vec![];
        for element in fragment.select(&selector) {
            let mut result = HashMap::new();

            let element = element.inner_html().replace(&target, "");

            let elements: Vec<&str> = element.split("-").collect();

            result.insert(String::from("name"), elements[0].trim().to_string());
            result.insert(String::from("type"), String::from("person"));
            result.insert(String::from("job_title"), elements[1].trim().to_string());
            result.insert(
                String::from("parent"),
                target_id.value.as_ref().unwrap().clone(),
            );

            results.push(Target::try_from(result)?);
        }

        Ok(results)
    }
}

impl LinkedinSearch {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(LinkedinSearch {})
    }
}

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
    Target, TargetType,
};

use crate::CompiledModule;

#[derive(Clone)]
pub struct EmailToDomain {}

#[async_trait]
impl CompiledModule for EmailToDomain {
    fn name(&self) -> String {
        "email_to_domain".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "search domain for specified email address".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    fn target_type(&self) -> TargetType {
        TargetType::Email
    }

    async fn run(
        &self,
        params: Args,
        _tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let email = params.get("target").unwrap();
        let target_id = params
            .get("target_id")
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let target = email.value.unwrap();

        let url = String::from("https://viewdns.info/reversewhois/?q=")
            .add(&target.replace(" ", "+").as_str());
        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/110.0"
                .parse()
                .map_err(|_| {
                    ErrorKind::Module(ModuleError::Execution("can-t set header".to_string()))
                })?,
        );
        headers.insert(
            "Accept",
            "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8"
                .parse()
                .map_err(|_| {
                    ErrorKind::Module(ModuleError::Execution("can-t set header".to_string()))
                })?,
        );
        headers.insert(
            "Accept-Language",
            "en-US,en;q=0.5".parse().map_err(|_| {
                ErrorKind::Module(ModuleError::Execution("can-t set header".to_string()))
            })?,
        );

        let client = reqwest::Client::builder()
            .default_headers(headers)
            .redirect(reqwest::redirect::Policy::none())
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
        let selector = Selector::parse("table[id=null] td table tr").map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(
                "invalid selector from scrapping".to_string(),
            ))
        })?;

        let mut results = vec![];

        let count = fragment.select(&selector).count();
        if count > 1 {
            for element in fragment.select(&selector) {
                if element.html().contains("Domain Name") {
                    continue;
                }

                let mut result = HashMap::new();
                let selector = Selector::parse("tr td").map_err(|_| {
                    ErrorKind::Module(ModuleError::Execution(
                        "invalid selector from scrapping".to_string(),
                    ))
                })?;

                let elements = element.select(&selector);
                let count = element.select(&selector).count();

                if count != 3 {
                    continue;
                }

                for (k, e) in elements.enumerate() {
                    if k == 0 {
                        result.insert(String::from("name"), e.inner_html().trim().to_string());
                        result.insert(String::from("type"), String::from("domain"));
                        result.insert(String::from("parent"), target_id.clone());
                    } else if k == 1 {
                        result.insert(
                            String::from("registered_at"),
                            e.inner_html().trim().to_string(),
                        );
                    }
                }
                if result.len() > 0 {
                    results.push(Target::try_from(result)?);
                }
            }
        }

        Ok(results)
    }
}

impl EmailToDomain {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(EmailToDomain {})
    }
}

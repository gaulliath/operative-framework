use std::collections::HashMap;

use async_trait::async_trait;

use reqwest::header;
use serde::Deserialize;
use serde::Serialize;
use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};

use crate::CompiledModule;

#[derive(Clone)]
pub struct AccountSearch {}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
struct Root {
    pub license: Vec<String>,
    pub authors: Vec<String>,
    pub categories: Vec<String>,
    pub sites: Vec<Site>,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Site {
    pub name: String,
    #[serde(rename = "uri_check")]
    pub uri_check: String,
    #[serde(rename = "e_code")]
    pub e_code: i64,
    #[serde(rename = "e_string")]
    pub e_string: String,
    #[serde(rename = "m_string")]
    pub m_string: String,
    #[serde(rename = "m_code")]
    pub m_code: i64,
    pub known: Vec<String>,
    pub cat: String,
    pub valid: bool,
}

#[async_trait]
impl CompiledModule for AccountSearch {
    fn name(&self) -> String {
        "accounts.search".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn is_threaded(&self) -> bool {
        true
    }

    fn resume(&self) -> String {
        "Search if username is used in multiples sources".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("target", false, false, None),
        ]
    }

    async fn run(
        &self,
        params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let company = params.get("target").unwrap();
        let target_id = params
            .get("target_id")
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let target = company
            .value
            .ok_or(ErrorKind::Module(ModuleError::TargetNotAvailable))?;
        let tx = tx.ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(
            "tx is mandatory for threaded module".to_string(),
        )))?;
        let websites = self.get_accounts().await?;
        for site in websites {
            let _ = tokio::spawn(AccountSearch::check_website(
                site,
                (target.clone(), target_id.clone()),
                tx.clone(),
            ));
        }
        Ok(vec![])
    }

    fn target_type(&self) -> TargetType {
        TargetType::Username
    }
}

impl AccountSearch {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(AccountSearch {})
    }

    pub async fn get_accounts(&self) -> Result<Vec<Site>, ErrorKind> {
        let url = "https://raw.githubusercontent.com/WebBreacher/WhatsMyName/main/wmn-data.json";

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

        let res: Root = client
            .get(url)
            .send()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?
            .json()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        Ok(res.sites)
    }

    async fn check_website(site: Site, account: (String, String), tx: UnboundedSender<Event>) {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            header::HeaderValue::from_static(
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
            ),
        );

        let (account, target_id) = account;

        if let Ok(client) = reqwest::Client::builder().default_headers(headers).build() {
            let url = site.uri_check.replace("{account}", &account);

            if let Ok(res) = client.get(&url).send().await {
                let status = res.status().to_string();
                if let Ok(content) = res.text().await {
                    if site.m_code != site.e_code {
                        if site.e_code.to_string().ne(status.as_str()) {
                            return;
                        }
                    }

                    if !content.contains(&site.e_string)
                        || (!site.m_string.is_empty() && content.contains(&site.m_string))
                    {
                        return;
                    }

                    let mut result = HashMap::new();
                    result.insert(String::from("name"), site.name.clone());
                    result.insert(String::from("type"), String::from("account"));
                    result.insert(String::from("category"), site.cat.clone());
                    result.insert(String::from("parent"), target_id.clone());
                    result.insert(String::from("link"), url);

                    if let Ok(target) = Target::try_from(result) {
                        let _ = send_event(&tx, Event::ResultsModule(vec![target])).await;
                        let _ = send_event(
                            &tx,
                            Event::ResponseSimple(format!(
                                "account {} found : {}",
                                account, site.name
                            )),
                        )
                        .await;
                    }
                }
            }
        }
    }
}

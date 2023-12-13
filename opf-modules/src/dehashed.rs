use async_trait::async_trait;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;

use reqwest::header;
use tokio::sync::mpsc::UnboundedSender;

use opf_models::error::{ErrorKind, Module as ModuleError};
use opf_models::event::{send_event, Event};
use opf_models::{
    metadata::{Arg, Args},
    Target, TargetType,
};

use crate::CompiledModule;

#[derive(Clone)]
pub struct DeHashed {}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Results {
    pub balance: i64,
    pub entries: Vec<Entry>,
    pub success: bool,
    pub took: String,
    pub total: i64,
}

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Entry {
    pub id: String,
    pub email: String,
    #[serde(rename = "ip_address")]
    pub ip_address: String,
    pub username: String,
    pub password: String,
    #[serde(rename = "hashed_password")]
    pub hashed_password: String,
    pub name: String,
    pub vin: String,
    pub address: String,
    pub phone: String,
    #[serde(rename = "database_name")]
    pub database_name: String,
}

#[async_trait]
impl CompiledModule for DeHashed {
    fn name(&self) -> String {
        "dehashed".to_string()
    }

    fn author(&self) -> String {
        "Tristan Granier".to_string()
    }

    fn resume(&self) -> String {
        "search information from data breach".to_string()
    }

    fn args(&self) -> Vec<Arg> {
        vec![
            Arg::new("target_id", true, false, None),
            Arg::new("type", false, true, Some("domain".to_string())),
            Arg::new("api_key", false, false, None),
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
        params: Args,
        tx: Option<UnboundedSender<Event>>,
    ) -> Result<Vec<Target>, ErrorKind> {
        let search_type =
            params
                .get("type")
                .ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(format!(
                    "search_type"
                ))))?;

        let api_key =
            params
                .get("api_key")
                .ok_or(ErrorKind::Module(ModuleError::ParamNotAvailable(format!(
                    "api_key"
                ))))?;

        let target_id = target.target_id.to_string();

        let target_value = target.target_name.clone();

        let search_type_value =
            search_type
                .value
                .ok_or(ErrorKind::Module(ModuleError::Execution(format!(
                    "can't parse value of argument 'search_type'"
                ))))?;

        let api_key_value = api_key
            .value
            .ok_or(ErrorKind::Module(ModuleError::Execution(format!(
                "can't parse value of argument 'api_key'"
            ))))?;

        if !api_key_value.contains(":") {
            return Err(ErrorKind::Module(ModuleError::Execution(format!(
                "format of api_key isn't valid please use : username:apikey"
            ))));
        }

        let creds = api_key_value.split(":").collect::<Vec<&str>>();

        let url = format!(
            "https://api.dehashed.com/search?query={}:{}",
            search_type_value, target_value
        );

        let mut headers = header::HeaderMap::new();
        headers.insert("Accept", "application/json".parse().unwrap());

        let client = reqwest::Client::builder()
            .default_headers(headers)
            .build()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let res = client
            .get(url)
            .basic_auth(creds[0], Some(creds[1]))
            .send()
            .await
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let res = res.json::<Results>().await.map_err(|_| {
            ErrorKind::Module(ModuleError::Execution(format!(
                "no results found for '{}'",
                target.target_name
            )))
        })?;

        let tx = tx.unwrap();
        let _ = send_event(
            &tx,
            Event::ResponseSimple(format!("available balance '{}'", res.balance)),
        )
        .await;

        let _ = send_event(
            &tx,
            Event::UpdateTargetMeta((
                target_id.clone(),
                (opf_models::email::IS_VALID.to_string(), "true".to_string()),
            )),
        )
        .await;

        let mut results = vec![];

        for entry in res.entries {
            let mut target = HashMap::new();
            target.insert(
                String::from(opf_models::target::NAME),
                entry.database_name.clone(),
            );
            target.insert(
                String::from(opf_models::breach::PASSWORD),
                entry.password.clone(),
            );
            target.insert(
                String::from(opf_models::breach::ENCRYPTED_PASSWORD),
                entry.hashed_password.clone(),
            );
            target.insert(
                String::from(opf_models::breach::USERNAME),
                entry.username.clone(),
            );
            target.insert(
                String::from(opf_models::target::TYPE),
                TargetType::Breach.to_string(),
            );
            target.insert(String::from(opf_models::target::PARENT), target_id.clone());
            results.push(Target::try_from(target)?);
        }

        Ok(results)
    }
}

impl DeHashed {
    pub fn new() -> Box<dyn CompiledModule> {
        Box::new(DeHashed {})
    }
}

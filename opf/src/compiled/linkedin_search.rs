use std::collections::HashMap;
use crate::compiled::Compiled;
use crate::core::session::Session;
use crate::error::{ErrorKind, Module as ModuleError};
use crate::modules::metadata::{Arg, Args};
use reqwest::header;
use std::ops::Add;
use scraper::{Html, Selector};
use crate::common::Target;

pub struct LinkedinSearch {}

impl Compiled for LinkedinSearch {
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
            Arg::new("company", true, false, None),
            Arg::new("limit", false, true, Some("10".to_string())),
        ]
    }

    fn run(&self, sess: &mut Session, params: Args) -> Result<Vec<Target>, ErrorKind> {
        println!("Execution {}...", &self.name());
        println!("Session {}", sess.session_id);
        println!("Params: {:#?}", params);

        let company = params.get("company").unwrap();
        let limit = params.get("limit").unwrap();
        let target_id = company.value.unwrap().parse::<i32>().unwrap();
        let target = sess.get_target_by_id(target_id).unwrap();

        let url = String::from("https://www.google.com/search?num=")
            .add(&limit.value.unwrap().as_str())
            .add("&start=0&hl=en&q=site:linkedin.com/in+")
            .add(&target.target_name.replace(" ", "+").as_str());

        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            header::HeaderValue::from_static(
                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
            ),
        );

        let client = reqwest::blocking::Client::builder()
            .default_headers(headers)
            .build()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let res = client
            .get(url)
            .send()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?
            .text()
            .map_err(|e| ErrorKind::Module(ModuleError::Execution(e.to_string())))?;

        let fragment = Html::parse_fragment(res.as_str());
        let selector = Selector::parse(".g h3")
            .map_err(|_| ErrorKind::Module(ModuleError::Execution("invalid selector from scrapping".to_string())))?;

        let mut results = vec![];
        for element in fragment.select(&selector) {
            let mut result = HashMap::new();

            let element = element.inner_html().replace(&target.target_name, "");

            let elements :  Vec<&str> =     element.split("-")
                .collect();

            result.insert(String::from("name"), elements[0].trim().to_string());
            result.insert(String::from("type"), String::from("person"));
            result.insert(String::from("job_title"), elements[1].trim().to_string());

            results.push(Target::try_from(result)?);
        }
        println!("results : {:?}", results);
        Ok(results)
    }
}

impl LinkedinSearch {
    pub fn new() -> Box<dyn Compiled> {
        Box::new(LinkedinSearch {})
    }
}

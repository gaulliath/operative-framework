use reqwest::header;
use rlua;
use std::collections::HashMap;

pub fn extends_http(ctx: &mut rlua::Context) {
    // extends here
    println!("extending module with http");

    let http_get = match ctx.create_function(|_, url: String| {
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
            .unwrap();

        let res = match client.get(url).send() {
            Ok(res) => res.text().unwrap(),
            Err(e) => {
                println!("ERR http_get {}", e);
                return Ok(None);
            }
        };
        Ok(Some(res))
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("http_get", http_get).unwrap();

    let http_get_json = match ctx.create_function(|_, url: String| {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            "User-Agent",
            header::HeaderValue::from_static(
                "Mozilla/5.0 (X11; Linux x86_64; rv:91.0) Gecko/20100101 Firefox/91.0",
            ),
        );

        let client = reqwest::blocking::Client::builder()
            .default_headers(headers)
            .build()
            .unwrap();

        let res: HashMap<String, String> = match client.get(url).send() {
            Ok(res) => res.json().unwrap(),
            Err(e) => {
                println!("ERR http_get_json {}", e);
                return Ok(None);
            }
        };
        Ok(Some(res))
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("functio http_get_json disabled ?");
            return;
        }
    };

    ctx.globals().set("http_get_json", http_get_json).unwrap()
}

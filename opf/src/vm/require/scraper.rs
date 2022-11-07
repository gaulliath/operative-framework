use scraper::{Html, Selector};

pub fn extends_scraper<'a>(ctx: &'a mut rlua::Context) {
    // extends here
    let scraper_get = match ctx.create_function(|_, (select, html): (String, String)| {

        let fragment = Html::parse_fragment(html.as_str());
        let selector = match Selector::parse(select.as_str()) {
            Ok(res) => res,
            Err(e) => {
                println!("ERR scraper {:?}", e);
                return Ok(vec![]);
            }
        };

        let mut results = vec![];
        for element in fragment.select(&selector) {
            results.push(element.inner_html())
        }
        Ok(results)

    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("scraper_get", scraper_get).unwrap();

}

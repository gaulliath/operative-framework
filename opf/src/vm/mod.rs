pub mod arg;
pub mod require;
mod target;
mod session;

use crate::core::session::Session;
use crate::modules::metadata::Requirements;

pub fn extends<'a, 'b, 'c>(
    ctx: &'a mut rlua::Context,
    require: &'b Requirements,
    session: &'c Session,
) {
    match require {
        Requirements::Http => require::http::extends_http(ctx),
        Requirements::Scraper => require::scraper::extends_scraper(ctx),
        Requirements::Target => require::target::extends_target(ctx),
        Requirements::Network => require::network::extends_network(ctx)
    }
}

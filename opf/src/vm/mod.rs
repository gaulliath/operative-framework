pub mod arg;
pub mod require;
mod session;
mod target;

use crate::core::session::Session;
use crate::modules::metadata::Requirements;

pub fn extends(ctx: &mut rlua::Context, require: &Requirements, _session: &Session) {
    match require {
        Requirements::Http => require::http::extends_http(ctx),
        Requirements::Scraper => require::scraper::extends_scraper(ctx),
        Requirements::Target => require::target::extends_target(ctx),
        Requirements::Network => require::network::extends_network(ctx),
        Requirements::Common => require::common::extends_common(ctx),
    }
}

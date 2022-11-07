use crate::core::session::Session;
use super::{Command, CommandAction};
use super::error::Error;

pub fn exec<'a>(session: &'a mut Session, cmd: Command) -> Result<(), Error> {
    match cmd.action {
        CommandAction::Save => save(session, cmd),
        CommandAction::Load => load(session, cmd),
        _ => Ok(()),
    }
}

/// save session into json file
fn save<'a>(session: &'a mut Session, cmd: Command) -> Result<(), Error> {
    let path = cmd.get("to");
    if path.is_none() {
        return Err(Error::GenericError(format!("argument '{}' is mandatory for session saving.", "to")));
    }

    let path = path.unwrap();
    println!("exporting session to {}", path);

    let open = std::fs::OpenOptions::new()
        .create_new(true)
        .write(true)
        .append(false)
        .open(path.clone());

    let file = match open {
        Ok(f) => f,
        Err(e) => return Err(Error::GenericError(e.to_string()))
    };

    match serde_json::to_writer(file, session) {
        Ok(json) => println!("session saved into {}", path),
        Err(e) => return Err(Error::GenericError(e.to_string()))
    };

    Ok(())
}

/// load json session file into current session
fn load<'a>(session: &'a mut Session, cmd: Command) -> Result<(), Error> {
    let path = cmd.get("from");
    if path.is_none() {
        return Err(Error::GenericError(format!("argument '{}' is mandatory for session loading.", "from")));
    }

    let path = path.unwrap();
    println!("loading session from {}", path);

    let file = match std::fs::File::open(path) {
        Ok(f) => f,
        Err(e) => return Err(Error::GenericError(e.to_string()))
    };

    let new_session: Session = match serde_json::from_reader(&file) {
        Ok(value) => value,
        Err(e) => return Err(Error::GenericError(e.to_string()))
    };

    *session = new_session;
    println!("session loaded successfully");
    Ok(())
}

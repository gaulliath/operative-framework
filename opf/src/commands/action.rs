use std::time::SystemTime;

use crate::commands::{Command, CommandAction};
use crate::core::session::Session;
use crate::error::ErrorKind;

pub fn exec(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    match cmd.action {
        CommandAction::Accept => accept(session, cmd),
        CommandAction::List => list(session, cmd),
        _ => Ok(()),
    }
}

fn list(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    let _params = cmd.params;
    let actions = &session.actions;
    let headers = vec![
        "id".to_string(),
        "type".to_string(),
        "state".to_string(),
        "by".to_string(),
    ];
    let mut rows = vec![];

    for action in actions {
        rows.push(vec![
            action.action_id.to_string(),
            action.action_type.to_string(),
            action.action_state.to_string(),
            action.action_created_by.clone(),
        ]);
    }
    session.output_table(headers, rows);
    Ok(())
}

fn accept(session: &mut Session, cmd: Command) -> Result<(), ErrorKind> {
    println!("accept action ! {:?}", cmd.params);
    let action = crate::common::Action {
        action_id: 1,
        action_type: crate::common::action::ActionType::MergeTarget,
        action_ctx: vec![],
        action_state: crate::common::action::ActionState::Pending,
        action_created_by: "linkedin.search".to_string(),
        action_created_at: SystemTime::now(),
    };

    session.actions.push(action);
    Ok(())
}

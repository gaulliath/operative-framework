use std::time::SystemTime;

use crate::commands::{Command, CommandAction};
use crate::common::action as opf_action;
use crate::core::session::Session;

pub fn exec<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_action::Error> {
    match cmd.action {
        CommandAction::Accept => accept(session, cmd),
        CommandAction::List => list(session, cmd),
        _ => Ok(()),
    }
}

fn list<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_action::Error> {
    let params = cmd.params;
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

fn accept<'a>(session: &'a mut Session, cmd: Command) -> Result<(), opf_action::Error> {
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

use crate::common::Target;
use crate::common::search::Target as TargetSearch;

pub fn extends_target<'a, 'b>(ctx: &'a mut rlua::Context) {
    // extends here
    println!("extending module with target");


    let get_target = match ctx.create_function(|this, target_id: String| {
        let targets : Vec<Target> = this.globals().get("sess_targets").unwrap();

        for target in &targets {
            if target.target_uuid.to_string() == target_id {
               return Ok(Some(target.clone()));
            } else if target.target_id.to_string() == target_id {
                return Ok(Some(target.clone()));
            }
        }
        Ok(None)
    }) {
        Ok(f) => f,
        Err(_) => {
            println!("function log disabled ?");
            return;
        }
    };

    ctx.globals().set("get_target", get_target).unwrap();
}

mod prompt;

use crate as opf;
use crossterm::event::{KeyCode, KeyModifiers};
use reedline::ExampleHighlighter;
use reedline::{ColumnarMenu, DefaultCompleter, ReedlineMenu, default_emacs_keybindings, Emacs};
use reedline::{ListMenu, Reedline, ReedlineEvent, Signal};

pub fn run(opf: &mut opf::core::Core) {
    //let (session, mut config) = opf::core::session::new();
    //config.verbose = false;

    //let mut opf = opf::core::new(session, config);
    //opf.init_manager("/home/graniet/.opf/modules/");

    let auto_complete = vec![
        "add target".into(),
        "list".into(),
        "list modules".into(),
        "list target".into(),
        "list group".into(),
    ];
    let completer = Box::new(DefaultCompleter::new_with_wordlen(auto_complete.clone(), 2));
    // Use the interactive menu to select options from the completer
    let completion_menu = Box::new(
        ColumnarMenu::default()
            .with_name("completion_menu")
            .with_only_buffer_difference(false),
    );

    let mut line_editor = Reedline::create()
        .with_completer(completer)
        .with_quick_completions(true)
        .with_partial_completions(true)
        .with_menu(ReedlineMenu::EngineCompleter(completion_menu))
        .with_edit_mode({
            let mut keybindings = default_emacs_keybindings();
            keybindings.add_binding(
                reedline::KeyModifiers::NONE,
                reedline::KeyCode::Tab,
                ReedlineEvent::UntilFound(vec![
                    ReedlineEvent::Menu("completion_menu".to_string()),
                    ReedlineEvent::MenuNext,
                ]),
            );

            Box::new(Emacs::new(keybindings))
        });
    //.with_highlighter(Box::new(ExampleHighlighter::new(commands)));
    let prompt = prompt::PromptCli {};
    loop {
        let sig = line_editor.read_line(&prompt).unwrap();
        match sig {
            Signal::Success(buffer) => {
                if let Err(e) = opf.run(buffer.trim().to_string()) {
                    println!("ERR {}", e);
                }
            }
            Signal::CtrlD | Signal::CtrlC => {
                return;
            }
        }
    }
}

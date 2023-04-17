use thiserror::Error;

#[derive(Error, Debug)]
pub enum ErrorKind {
    #[error("argument is not found for this command")]
    InvalidCommandArgument,
    #[error("invalid format for argument")]
    InvalidFormatArgument,
    #[error("action not available")]
    ActionNotAvailable,
    #[error("context object for this action not available")]
    ObjectNotAvailable,
    #[error("{0}")]
    GenericError(String),
    #[error(transparent)]
    Manager(Manager),
    #[error(transparent)]
    Action(Action),
    #[error(transparent)]
    Group(Group),
    #[error(transparent)]
    Link(Link),
    #[error(transparent)]
    Workspace(Workspace),
    #[error(transparent)]
    Module(Module),
    #[error(transparent)]
    Target(Target),
    #[error(transparent)]
    Export(Export),
    #[error(transparent)]
    Metadata(Metadata),
    #[error("channel communication error : {0}")]
    Channel(String),
}

#[derive(Error, Debug)]
pub enum Metadata {
    #[error("field required : {0}")]
    Required(String),
    #[error("field format incorrect")]
    FormatField,
    #[error("extender requirement not available : {0}")]
    Requirement(String),
}

#[derive(Error, Debug)]
pub enum Manager {
    #[error("can't process file")]
    CantProcessFile,
    #[error("directory {0} cant be opened")]
    CantOpenDirectory(String),
    #[error("can't read a content of module {0}")]
    CantReadContent(String),
    #[error("parsing metadata of {0} is impossible, please check format")]
    CantParseMetadata(String),
}

#[derive(Error, Debug)]
pub enum Action {
    #[error("action {0} is not available")]
    ActionNotFound(String),
}

#[derive(Error, Debug)]
pub enum Group {
    #[error("{0}")]
    GenericError(String),
    #[error("link already found in this session")]
    LinkExist,
    #[error("group can't retrieve group {0}")]
    Get(String),
    #[error("group already exist {0}")]
    Exist(String),
}

#[derive(Error, Debug)]
pub enum Workspace {
    #[error("{0}")]
    GenericError(String),
    #[error("group can't retrieve group {0}")]
    Get(String),
    #[error("workspace already exist {0}")]
    Exist(String),
    #[error("parameter '{0}' mandatory for this action")]
    ParamNotFound(String),
}

#[derive(Error, Debug)]
pub enum Link {
    #[error("parameter '{0}' mandatory for this action")]
    ParamNotFound(String),
    #[error("format invalid for {0} parameter")]
    ParamFormatInvalid(String),
    #[error("Parameter 'label' mandatory for this action")]
    ParamLabelNotFound,
    #[error("Parameter 'from' mandatory for this action")]
    ParamFromNotFound,
    #[error("Parameter 'to' mandatory for this action")]
    ParamToNotFound,
    #[error("Parameter 'type' mandatory for this action")]
    ParamTypeNotFound,
    #[error("link already found in this session")]
    LinkExist,
    #[error("link not found")]
    NotFound,
    #[error("target do not exist in session?")]
    TargetNotFound,
}

#[derive(Error, Debug)]
pub enum Module {
    #[error("module not found")]
    ModuleNameNotFound,
    #[error("group not found")]
    GroupNotFound,
    #[error("cant load content of module")]
    CantLoadContent,
    #[error("parameter {0} not specified, please set")]
    ParamNotAvailable(String),
    #[error("target not available in current session")]
    TargetNotAvailable,
    #[error("execution error {0}")]
    Execution(String),
}

#[derive(Error, Debug)]
pub enum Target {
    #[error("Parameter 'name' mandatory for this action")]
    ParamNameNotFound,
    #[error("Parameter 'type' mandatory for this action")]
    ParamTypeNotFound,
    #[error("Target already found in this session")]
    TargetExist,
    #[error("Target not found in this session")]
    NotFound,
    #[error("Target type '{0}' not available")]
    TypeNotAvailable(String),
    #[error("Parent not found")]
    ParentUuidNotFound,
    #[error("Parent uuid ins't valid")]
    ParentUuidNotValid,
    #[error("Target can't be created")]
    CantBeCreated,
    #[error("{0}")]
    Parsing(String),
}

#[derive(Error, Debug)]
pub enum Export {
    #[error("module not found")]
    ModuleNameNotFound,
    #[error("cant load content of module")]
    CantLoadContent,
    #[error("export type not found")]
    ExportType,
}

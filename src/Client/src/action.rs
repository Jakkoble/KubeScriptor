use crate::app::Job;

pub(crate) enum Action {
    Quit,
    SelectJob(Job),
    OpenJobList,
}

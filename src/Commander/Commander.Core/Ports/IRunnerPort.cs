using Commander.Core.Entities;

namespace Commander.Core.Ports;

public interface IRunnerPort
{
  Task ExecuteJob(Job job);
  Task StopJob(Job job, bool wasSuccessful);
}

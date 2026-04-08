using Commander.Core.Entities;

namespace Commander.Server.Store;

public interface IJobStore
{
  bool StoreJob(Job job);
  Job GetJob(Guid jobId);
}

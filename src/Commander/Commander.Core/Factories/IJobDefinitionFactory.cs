using Commander.Core.Entities;

namespace Commander.Core.Factories;

public interface IJobDefinitionFactory
{
  Job CreateFromYaml(string? yamlPayload);
}

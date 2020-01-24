import { noiseProvider } from './NoiseScenario';
import { waveProvider } from './WaveScenario';
import { arrowFileProvider } from './ArrowFileScenario';
import { csvWaveProvider } from './CSVWaveScenario';
import { ScenarioProvider, Scenario } from 'types';

export const scenarios: Record<Scenario, ScenarioProvider<any>> = {
  noise: noiseProvider,
  wave: waveProvider,
  arrowFile: arrowFileProvider,
  csvWave: csvWaveProvider,
};

export const scenarioList = Object.values(scenarios);

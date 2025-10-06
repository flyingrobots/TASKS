export interface Config {
  server: {
    port: number;
    host: string;
  };
  paths: {
    dagFile: string;
    featuresFile: string;
    wavesFile: string;
    tasksFile: string;
    stateFile: string;
    logsFile: string;
    dagStateServer: string;
  };
  vite: {
    port: number;
    host: string;
  };
}

const defaultConfig: Config = {
  server: {
    port: 8080,
    host: 'localhost'
  },
  paths: {
    dagFile: '',
    featuresFile: '',
    wavesFile: '',
    tasksFile: '',
    stateFile: '',
    logsFile: '',
    dagStateServer: ''
  },
  vite: {
    port: 5173,
    host: 'localhost'
  }
};

let config: Config = defaultConfig;

export async function loadConfig(): Promise<Config> {
  try {
    const response = await fetch('/config.json');
    if (response.ok) {
      config = await response.json();
      console.log('Loaded configuration:', config);
    } else {
      console.warn('Config file not found, using defaults');
    }
  } catch (error) {
    console.error('Error loading config:', error);
  }
  return config;
}

export function getConfig(): Config {
  return config;
}

export function getWebSocketUrl(): string {
  return `ws://${config.server.host}:${config.server.port}`;
}
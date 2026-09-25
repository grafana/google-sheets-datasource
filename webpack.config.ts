import { Configuration } from 'webpack';
import CopyWebpackPlugin from 'copy-webpack-plugin';
import { mergeWithRules } from 'webpack-merge';
import path from 'path';
import grafanaConfig from './.config/webpack/webpack.config';

const config = async (env: any): Promise<Configuration> => {
  const baseConfig = await grafanaConfig(env);
  const customConfig = {
    plugins: [
      new CopyWebpackPlugin({
        patterns: [
          {
            from: '../pkg/schema/dsconfig.json',
            to: './schema/dsconfig.json',
            noErrorOnMissing: true,
          },
          {
            from: '../pkg/schema/schema.gen.json',
            to: './schema/v0alpha1.json',
            noErrorOnMissing: true,
          },
          {
            from: '../pkg/schema/settings.gen.json',
            to: './schema/v0alpha1/settings.json',
            noErrorOnMissing: true,
          },
          {
            from: '../pkg/schema/settings.examples.gen.json',
            to: './schema/v0alpha1/settings.examples.json',
            noErrorOnMissing: true,
          },
        ],
      }),
    ],
  };

  return mergeWithRules({
    module: {
      rules: {
        exclude: 'replace',
      },
    },
  })(baseConfig, customConfig);
};

export default config;

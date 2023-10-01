import visualizer from 'rollup-plugin-visualizer'
import { ANALYSIS } from '../../constant'

export function ConfigVisualizerConfig(): any {
  if (ANALYSIS) {
    return visualizer({
      filename: 'dist/report.html',
      open: true,
      gzipSize: true,
      emitFile: false,
    })
  }
  return null
}

import { useState } from "react";
import "./App.css";
import {
  Chart as ChartJS,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Legend,
  Tooltip
} from "chart.js";
import { Line } from "react-chartjs-2";
import { Analyze, AnalyzeLib, ApplyFilters, ExportFilterResults } from "../wailsjs/go/main/App";

ChartJS.register(
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Legend,
  Tooltip
);

type SignalDTO = {
  samples: number[];
  sampleRate: number;
};

type SpectrumDTO = {
  freqs: number[];
  mag: number[];
  phase: number[];
};

type AnalysisResult = {
  x: SignalDTO;
  y: SignalDTO;
  conv: SignalDTO;
  corr: SignalDTO;
  spectrumX: SpectrumDTO;
  spectrumY: SpectrumDTO;
  spectrumConv: SpectrumDTO;
  spectrumCorr?: SpectrumDTO;
  idftX?: SignalDTO;
  idftY?: SignalDTO;
  idftConv?: SignalDTO;
};

type FilterResult = {
  filtered: SignalDTO;
  freqs: number[];
  magnitude: number[];
  inputSignal: SignalDTO;
  inputSpectrum: SpectrumDTO;
};

type FilterResults = {
  [key: string]: FilterResult;
};

function App() {
  const [data, setData] = useState<AnalysisResult | null>(null);
  const [dataLib, setDataLib] = useState<AnalysisResult | null>(null);
  const [filterResults, setFilterResults] = useState<FilterResults | null>(null);
  const [exportResults, setExportResults] = useState<{ [key: string]: string } | null>(null);
  const [N, setN] = useState<number>(256);
  const [audioContext, setAudioContext] = useState<AudioContext | null>(null);

  const initAudio = () => {
    if (!audioContext) {
      const ctx = new (window.AudioContext || (window as any).webkitAudioContext)();
      setAudioContext(ctx);
      return ctx;
    }
    return audioContext;
  };

  const playSignal = async (sig: SignalDTO, title: string) => {
    try {
      const ctx = initAudio();
      if (ctx.state === 'suspended') {
        await ctx.resume();
      }

      let maxVal = Math.max(...sig.samples.map(Math.abs));
      if (maxVal === 0) maxVal = 1;

      const normalized = sig.samples.map(s => s / maxVal * 0.5);

      const buffer = ctx.createBuffer(1, normalized.length, ctx.sampleRate);
      const channelData = buffer.getChannelData(0);

      for (let i = 0; i < normalized.length; i++) {
        channelData[i] = normalized[i];
      }

      const source = ctx.createBufferSource();
      source.buffer = buffer;
      source.connect(ctx.destination);
      source.start();

      console.log(`Воспроизведение: ${title}, sampleRate: ${ctx.sampleRate}, length: ${normalized.length}`);
    } catch (e) {
      console.error("Ошибка воспроизведения:", e);
      alert("Ошибка воспроизведения аудио. Разрешите автовоспроизведение в браузере.");
    }
  };

  const handleAnalyze = async () => {
    try {
      // @ts-ignore
      const res = (await Analyze(N)) as AnalysisResult;
      setData(res);

      if (res.x && res.y) {
        setTimeout(() => {
          playSignal(res.x, "x(t)").catch(console.error);
        }, 500);
        setTimeout(() => {
          playSignal(res.y, "y(t)").catch(console.error);
        }, 3000);
      }
    } catch (e) {
      console.error(e);
      alert("Ошибка при вызове Analyze() из Go. Смотри консоль.");
    }
  };

  const handleAnalyzeLib = async () => {
    try {
      // @ts-ignore
      const res = (await AnalyzeLib(N)) as AnalysisResult;
      setDataLib(res);

      if (res.x && res.y) {
        setTimeout(() => {
          playSignal(res.x, "x(t) библиотека").catch(console.error);
        }, 500);
        setTimeout(() => {
          playSignal(res.y, "y(t) библиотека").catch(console.error);
        }, 3000);
      }
    } catch (e) {
      console.error(e);
      alert("Ошибка при вызове AnalyzeLib() из Go. Смотри консоль.");
    }
  };

  const handleApplyFilters = async () => {
    try {
      // @ts-ignore
      const res = (await ApplyFilters(N)) as FilterResults;
      setFilterResults(res);

      if (res.movingAverage) {
        setTimeout(() => {
          playSignal(res.movingAverage.inputSignal, "Входной сигнал с помехами").catch(console.error);
        }, 500);
        setTimeout(() => {
          playSignal(res.movingAverage.filtered, "Однородный фильтр").catch(console.error);
        }, 3000);
        setTimeout(() => {
          playSignal(res.firBandStop.filtered, "КИХ режекторный фильтр").catch(console.error);
        }, 5500);
        setTimeout(() => {
          playSignal(res.iirBandPass.filtered, "БИХ полосовой фильтр").catch(console.error);
        }, 8000);
      }
    } catch (e) {
      console.error(e);
      alert("Ошибка при применении фильтров. Смотри консоль.");
    }
  };

  const handleExportFilters = async () => {
    try {
      // @ts-ignore
      const res = (await ExportFilterResults(N)) as { [key: string]: string };
      setExportResults(res);
      alert("Аудиофайлы успешно сохранены в папку 'output'!");
    } catch (e) {
      console.error(e);
      alert("Ошибка при экспорте фильтров. Смотри консоль.");
    }
  };

  const renderSignalChart = (sig?: SignalDTO, title?: string) => {
    if (!sig) return null;
    const N = sig.samples.length;
    const labels = Array.from({ length: N }, (_, i) => (i / sig.sampleRate).toFixed(5));
    return (
      <div style={{ marginBottom: 24 }}>
        <h3 style={{ color: "#7b2c34" }}>{title}</h3>
        <Line
          data={{
            labels,
            datasets: [
              {
                label: title,
                data: sig.samples,
                borderColor: "#a23e48",
                backgroundColor: "rgba(162, 62, 72, 0.08)",
                borderWidth: 2,
                pointRadius: 0,
                fill: false
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: {
                title: { display: true, text: "t, s", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: {
                  maxTicksLimit: 10,
                  maxRotation: 0,
                  autoSkip: true,
                  color: "#7a5c4d"
                }
              },
              y: {
                title: { display: true, text: "Amplitude", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: { color: "#7a5c4d" }
              }
            }
          }}
        />
      </div>
    );
  };

  const renderSpectrumChart = (spec?: SpectrumDTO, title?: string) => {
    if (!spec) return null;
    return (
      <div style={{ marginBottom: 24 }}>
        <h3 style={{ color: "#7b2c34" }}>{title} — амплитудный спектр</h3>
        <Line
          data={{
            labels: spec.freqs,
            datasets: [
              {
                label: "Magnitude",
                data: spec.mag,
                borderColor: "#7b2c34",
                backgroundColor: "rgba(123, 44, 52, 0.08)",
                borderWidth: 2,
                pointRadius: 0,
                fill: false
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: {
                title: { display: true, text: "f, Hz", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: {
                  maxTicksLimit: 20,
                  maxRotation: 0,
                  autoSkip: true,
                  color: "#7a5c4d"
                }
              },
              y: {
                title: { display: true, text: "|X(f)|", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: { color: "#7a5c4d" }
              }
            }
          }}
        />
        <h4 style={{ color: "#8b5cf6" }}>{title} — фазовый спектр</h4>
        <Line
          data={{
            labels: spec.freqs,
            datasets: [
              {
                label: "Phase",
                data: spec.phase,
                borderColor: "#8b5cf6",
                backgroundColor: "rgba(139, 92, 246, 0.08)",
                borderWidth: 2,
                pointRadius: 0,
                fill: false
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: {
                title: { display: true, text: "f, Hz", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: { color: "#7a5c4d" }
              },
              y: {
                title: { display: true, text: "Phase, rad", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: { color: "#7a5c4d" }
              }
            }
          }}
        />
      </div>
    );
  };

  const renderFilterFrequencyResponse = (freqs?: number[], mag?: number[], title?: string) => {
    if (!freqs || !mag || freqs.length === 0) return null;
    return (
      <div style={{ marginBottom: 24 }}>
        <h3 style={{ color: "#7b2c34" }}>{title} — АЧХ</h3>
        <Line
          data={{
            labels: freqs,
            datasets: [
              {
                label: "Magnitude Response",
                data: mag,
                borderColor: "#059669",
                backgroundColor: "rgba(5, 150, 105, 0.08)",
                borderWidth: 2,
                pointRadius: 0,
                fill: false
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: {
                title: { display: true, text: "f, Hz", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: {
                  maxTicksLimit: 20,
                  maxRotation: 0,
                  autoSkip: true,
                  color: "#7a5c4d"
                }
              },
              y: {
                title: { display: true, text: "|H(f)|", color: "#7a5c4d" },
                grid: { color: "#e6d3c3" },
                ticks: { color: "#7a5c4d" }
              }
            }
          }}
        />
      </div>
    );
  };

  return (
    <div style={{ padding: 24, fontFamily: "system-ui" }}>
      <h2 style={{ color: "#7b2c34" }}>1 лабораторная</h2>
      <p style={{ color: "#7b2c34" }}>
        A = [0.8, 0.5, 0.3], f0 = 330 Гц, h = [1, 2, 3], φx = 0
      </p>
      <p style={{ color: "#7b2c34" }}>
        A = [0.8, 0.5, 0.3], f0 = 330 Гц, h = [1, 2, 3] φy = π/2.
      </p>

      <div style={{ marginBottom: 24, display: "flex", gap: 12, alignItems: "center", justifyContent: "center", flexWrap: "wrap" }}>
        <input
          type="number"
          value={N}
          onChange={(e) => setN(parseInt(e.target.value) || 256)}
          min="64"
          max="4096"
          style={{ marginLeft: 8, padding: 4, width: 100, color: "7b2c34" }}
        />
        <button
          onClick={handleAnalyze}
          style={{
            background: "#7b2c34",
            color: "#fff",
            border: "none",
            padding: "8px 14px",
            borderRadius: 6,
            cursor: "pointer"
          }}
        >
          Выполнить анализ (наше)
        </button>
        <button
          onClick={handleAnalyzeLib}
          style={{
            background: "#a23e48",
            color: "#fff",
            border: "none",
            padding: "8px 14px",
            borderRadius: 6,
            cursor: "pointer"
          }}
        >
          Выполнить анализ (библиотека gonum)
        </button>
      </div>

      {data && (
        <>
          {renderSignalChart(data.x, "x(t)")}
          {renderSignalChart(data.y, "y(t)")}
          {renderSpectrumChart(data.spectrumX, "ДПФ X(f)")}
          {data.idftX && renderSignalChart(data.idftX, "ОДПФ x(t)")}
          {renderSpectrumChart(data.spectrumX, "БПФ X(f)")}
          {data.idftX && renderSignalChart(data.idftX, "ОБПФ x(t)")}
          {renderSpectrumChart(data.spectrumY, "ДПФ Y(f)")}
          {data.idftY && renderSignalChart(data.idftY, "ОДПФ y(t)")}
          {renderSpectrumChart(data.spectrumY, "БПФ Y(f)")}
          {data.idftY && renderSignalChart(data.idftY, "ОБПФ y(t)")}
          {renderSignalChart(data.conv, "Свертка")}
          {dataLib && renderSignalChart(dataLib.conv, "Свертка через БПФ")}
          {renderSignalChart(data.corr, "Корреляция")}
          {dataLib && renderSignalChart(dataLib.corr, "Корреляция через БПФ")}
        </>
      )}

      {dataLib && (
        <>
          {renderSpectrumChart(dataLib.spectrumX, "БПФ X(f) (OS)")}
          {renderSpectrumChart(dataLib.spectrumY, "БПФ Y(f) (OS)")}
          {renderSignalChart(dataLib.conv, "Свертка (OS)")}
          {renderSignalChart(dataLib.corr, "Корреляция (OS)")}
        </>
      )}

      <hr style={{ borderColor: "#7b2c34", margin: "48px 0" }} />

      <h2 style={{ color: "#7b2c34" }}>2 лабораторная - Фильтры</h2>
      <p style={{ color: "#7b2c34" }}>
        Входной сигнал: 330 Гц (инструмент), фаза=0<br />
        Помеха: 660 Гц<br />
        Однородный фильтр: M=15<br />
        КИХ режекторный фильтр (Ханна): [653,667] Гц, M=201<br />
        БИХ полосовой фильтр: f0=400 Гц, BW=80 Гц
      </p>

      <div style={{ marginBottom: 24, display: "flex", gap: 12 }}>
        <button
          onClick={handleApplyFilters}
          style={{
            background: "#059669",
            color: "#fff",
            border: "none",
            padding: "8px 14px",
            borderRadius: 6,
            cursor: "pointer"
          }}
        >
          Применить фильтры
        </button>
        <button
          onClick={handleExportFilters}
          style={{
            background: "#7b5f47",
            color: "#fff",
            border: "none",
            padding: "8px 14px",
            borderRadius: 6,
            cursor: "pointer"
          }}
        >
          Экспортировать в WAV
        </button>
      </div>

      {filterResults && (
        <>
          <h3 style={{ color: "#7b2c34" }}>Входной сигнал с помехами</h3>
          {renderSignalChart(filterResults.movingAverage.inputSignal, "Входной сигнал")}
          {renderSpectrumChart(filterResults.movingAverage.inputSpectrum, "Спектр входного сигнала")}

          <h3 style={{ color: "#7b2c34" }}>1. Однородный фильтр (M=15)</h3>
          {renderSignalChart(filterResults.movingAverage.filtered, "Выход однородного фильтра")}

          <h3 style={{ color: "#7b2c34" }}>2. КИХ режекторный фильтр (653-667 Гц, M=201, Ханна)</h3>
          {renderFilterFrequencyResponse(filterResults.firBandStop.freqs, filterResults.firBandStop.magnitude, "КИХ фильтр")}
          {renderSignalChart(filterResults.firBandStop.filtered, "Выход КИХ режекторного фильтра")}

          <h3 style={{ color: "#7b2c34" }}>3. БИХ полосовой фильтр (f0=400 Гц, BW=80 Гц)</h3>
          {renderFilterFrequencyResponse(filterResults.iirBandPass.freqs, filterResults.iirBandPass.magnitude, "БИХ фильтр")}
          {renderSignalChart(filterResults.iirBandPass.filtered, "Выход БИХ полосового фильтра")}
        </>
      )}

      {exportResults && (
        <div style={{ marginTop: 48, padding: 16, backgroundColor: "#f5f1ed", borderRadius: 8 }}>
          <h3 style={{ color: "#7b2c34" }}>Экспортированные файлы</h3>
          <ul style={{ color: "#7b2c34", lineHeight: 2 }}>
            {exportResults.input && <li><strong>Входной сигнал:</strong> {exportResults.input}</li>}
            {exportResults.moving_average && <li><strong>Однородный фильтр:</strong> {exportResults.moving_average}</li>}
            {exportResults.fir_bandstop && <li><strong>КИХ режекторный фильтр:</strong> {exportResults.fir_bandstop}</li>}
            {exportResults.iir_bandpass && <li><strong>БИХ полосовой фильтр:</strong> {exportResults.iir_bandpass}</li>}
          </ul>
          <p style={{ color: "#666", fontSize: 14 }}>Файлы сохранены в папку 'output' в корне проекта</p>
        </div>
      )}
    </div>
  );
}

export default App;

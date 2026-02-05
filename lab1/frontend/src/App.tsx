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
import { Analyze, AnalyzeLib } from "../wailsjs/go/main/App";

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

function App() {
  const [data, setData] = useState<AnalysisResult | null>(null);
  const [dataLib, setDataLib] = useState<AnalysisResult | null>(null);
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

    const timeLabels = (sig: SignalDTO) => {
        const N = sig.samples.length;
        const dt = 1 / sig.sampleRate;
        return Array.from({ length: N }, (_, i) => +(i * dt).toFixed(6)); // немного округлим
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

  return (
    <div style={{ padding: 24, fontFamily: "system-ui" }}>
      <h2 style={{ color: "#7b2c34" }}>1 лабараторная</h2>
      <p style={{ color: "#7b2c34" }}>
        A = [0.8, 0.5, 0.3], f0 = 330 Гц, h = [1, 2, 3], φx = 0
      </p>
        <p style={{ color: "#7b2c34" }}>
            A = [0.8, 0.5, 0.3], f0 = 330 Гц, h = [1, 2, 3] φy = π/2.
        </p>

      <div style={{ marginBottom: 24, display: "flex", gap: 12, alignItems: "center", justifyContent: "center" }}>
          <input
            type="number"
            value={N}
            onChange={(e) => setN(parseInt(e.target.value) || 256)}
            min="64"
            max="4096"
            style={{ marginLeft: 8, padding: 4, width: 100, color: "7b2c34"}}
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
    </div>
  );
}

export default App;

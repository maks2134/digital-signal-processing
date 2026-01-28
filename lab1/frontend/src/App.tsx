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
import { Analyze } from "../wailsjs/go/main/App";

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
};

function App() {
  const [data, setData] = useState<AnalysisResult | null>(null);

  const handleAnalyze = async () => {
    try {
      const res = (await Analyze()) as AnalysisResult;
      setData(res);
    } catch (e) {
      console.error(e);
      alert("Ошибка при вызове Analyze() из Go. Смотри консоль.");
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
    const labels = Array.from({ length: N }, (_, i) => i / sig.sampleRate);
    return (
      <div style={{ marginBottom: 24 }}>
        <h3>{title}</h3>
        <Line
          data={{
            labels,
            datasets: [
              {
                label: title,
                data: sig.samples,
                borderColor: "rgba(75, 192, 192, 1)",
                borderWidth: 1,
                pointRadius: 0
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: { title: { display: true, text: "t, s" } },
              y: { title: { display: true, text: "Amplitude" } }
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
        <h3>{title} — амплитудный спектр</h3>
        <Line
          data={{
            labels: spec.freqs,
            datasets: [
              {
                label: "Magnitude",
                data: spec.mag,
                borderColor: "rgba(255, 99, 132, 1)",
                borderWidth: 1,
                pointRadius: 0
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: { title: { display: true, text: "f, Hz" } },
              y: { title: { display: true, text: "|X(f)|" } }
            }
          }}
        />
        <h4>{title} — фазовый спектр</h4>
        <Line
          data={{
            labels: spec.freqs,
            datasets: [
              {
                label: "Phase",
                data: spec.phase,
                borderColor: "rgba(54, 162, 235, 1)",
                borderWidth: 1,
                pointRadius: 0
              }
            ]
          }}
          options={{
            responsive: true,
            scales: {
              x: { title: { display: true, text: "f, Hz" } },
              y: { title: { display: true, text: "Phase, rad" } }
            }
          }}
        />
      </div>
    );
  };

  return (
    <div style={{ padding: 24, fontFamily: "system-ui" }}>
      <h2>Лаба 1 — свертка, корреляция, ДПФ</h2>
      <p>
        Параметры варианта: A = [0.8, 0.5, 0.3], f0 = 330 Гц, h = [1, 2, 3], φx
        = 0, φy = π/2.
      </p>
      <button onClick={handleAnalyze} style={{ marginBottom: 24 }}>
        Выполнить анализ (через Go/Wails)
      </button>

      {data && (
        <>
          {renderSignalChart(data.x, "x(t)")}
          {renderSignalChart(data.y, "y(t)")}
          {renderSignalChart(data.conv, "Свертка x*y")}
          {renderSignalChart(data.corr, "Корреляция Rxy")}

          {renderSpectrumChart(data.spectrumX, "X(f)")}
          {renderSpectrumChart(data.spectrumY, "Y(f)")}
          {renderSpectrumChart(data.spectrumConv, "Conv(f)")}
        </>
      )}
    </div>
  );
}

export default App;

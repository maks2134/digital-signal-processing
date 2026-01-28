export namespace internal {
	
	export class SpectrumDTO {
	    freqs: number[];
	    mag: number[];
	    phase: number[];
	
	    static createFrom(source: any = {}) {
	        return new SpectrumDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.freqs = source["freqs"];
	        this.mag = source["mag"];
	        this.phase = source["phase"];
	    }
	}
	export class SignalDTO {
	    samples: number[];
	    sampleRate: number;
	
	    static createFrom(source: any = {}) {
	        return new SignalDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.samples = source["samples"];
	        this.sampleRate = source["sampleRate"];
	    }
	}
	export class AnalysisResult {
	    x: SignalDTO;
	    y: SignalDTO;
	    conv: SignalDTO;
	    corr: SignalDTO;
	    spectrumX: SpectrumDTO;
	    spectrumY: SpectrumDTO;
	    spectrumConv: SpectrumDTO;
	
	    static createFrom(source: any = {}) {
	        return new AnalysisResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.x = this.convertValues(source["x"], SignalDTO);
	        this.y = this.convertValues(source["y"], SignalDTO);
	        this.conv = this.convertValues(source["conv"], SignalDTO);
	        this.corr = this.convertValues(source["corr"], SignalDTO);
	        this.spectrumX = this.convertValues(source["spectrumX"], SpectrumDTO);
	        this.spectrumY = this.convertValues(source["spectrumY"], SpectrumDTO);
	        this.spectrumConv = this.convertValues(source["spectrumConv"], SpectrumDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}


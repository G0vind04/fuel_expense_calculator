export namespace main {
	
	export class FuelCalculation {
	    distance: number;
	    fuelEfficiency: number;
	    fuelPrice: number;
	    fuelNeeded: number;
	    totalCost: number;
	    costPerKm: number;
	
	    static createFrom(source: any = {}) {
	        return new FuelCalculation(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.distance = source["distance"];
	        this.fuelEfficiency = source["fuelEfficiency"];
	        this.fuelPrice = source["fuelPrice"];
	        this.fuelNeeded = source["fuelNeeded"];
	        this.totalCost = source["totalCost"];
	        this.costPerKm = source["costPerKm"];
	    }
	}
	export class TripComparison {
	    vehicle1: FuelCalculation;
	    vehicle2: FuelCalculation;
	    savings: number;
	
	    static createFrom(source: any = {}) {
	        return new TripComparison(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.vehicle1 = this.convertValues(source["vehicle1"], FuelCalculation);
	        this.vehicle2 = this.convertValues(source["vehicle2"], FuelCalculation);
	        this.savings = source["savings"];
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


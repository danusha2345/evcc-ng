import { describe, expect, test } from "vitest";
import { fmtEnergy, optionStep } from "./energyOptions";
import { POWER_UNIT } from "../mixins/formatter";

// minimal fmtWh stub mirroring Intl number formatting with fixed fraction digits
const fmtWh = (wh: number, unit: POWER_UNIT, _withUnit: boolean, digits?: number) => {
  const inKWh = unit === POWER_UNIT.KW;
  const value = inKWh ? wh / 1e3 : wh;
  return `${value.toFixed(digits ?? 0)} ${inKWh ? "kWh" : "Wh"}`;
};

describe("fmtEnergy", () => {
  // #30736: a fractional limit on a coarse integer step (e.g. step 5 for a
  // no-battery device, maxEnergy 100) must keep its decimal instead of rounding
  test("shows a decimal for a fractional value even with an integer step", () => {
    expect(fmtEnergy(2.4, optionStep(100), fmtWh, "off")).toBe("2.4 kWh"); // step 5
  });

  test("keeps clean integer values without a decimal", () => {
    expect(fmtEnergy(5, optionStep(100), fmtWh, "off")).toBe("5 kWh");
    expect(fmtEnergy(10, optionStep(25), fmtWh, "off")).toBe("10 kWh"); // step 1
  });

  test("zero returns the zero text", () => {
    expect(fmtEnergy(0, 5, fmtWh, "off")).toBe("off");
  });

  test("fractional step keeps one decimal as before", () => {
    expect(fmtEnergy(1, optionStep(8), fmtWh, "off")).toBe("1.0 kWh"); // step 0.5
    expect(fmtEnergy(2.5, optionStep(8), fmtWh, "off")).toBe("2.5 kWh");
  });
});

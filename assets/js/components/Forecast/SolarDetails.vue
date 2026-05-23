<template>
	<div v-if="days.length">
		<div v-if="currentPower !== null" class="mb-2">
			<small>
				<span class="text-gray">{{ $t("forecast.solar.now") }}</span>
				<span class="text-primary fw-bold ms-2">{{
					fmtW(currentPower, POWER_UNIT.AUTO)
				}}</span>
			</small>
		</div>
		<div class="row gx-2 mt-1">
			<div v-for="day in days" :key="day.key" class="col-4" :class="`text-${day.align}`">
				<small>
					<span class="text-gray">{{ day.label }}</span>
					<br />
					<div
						class="d-flex flex-column flex-md-row column-gap-2"
						:class="`justify-content-md-${day.align}`"
					>
						<span class="text-primary fw-bold">{{ day.energy }}</span>
						<span v-if="day.note" class="text-gray">{{ day.note }}</span>
					</div>
				</small>
			</div>
		</div>
	</div>
</template>

<script lang="ts">
import { defineComponent, type PropType } from "vue";
import formatter, { POWER_UNIT } from "@/mixins/formatter";
import type { SolarDetails } from "./types";

export default defineComponent({
	name: "SolarDetails",
	mixins: [formatter],
	props: {
		solar: { type: Object as PropType<SolarDetails> },
	},
	computed: {
		// currentPower returns the forecasted solar power for the current
		// time slot (evcc-io/evcc#19649), or null if no timeseries is
		// available. timeseries is ordered; we take the last slot whose
		// timestamp is at or before now.
		currentPower(): number | null {
			const ts = this.solar?.timeseries;
			if (!ts?.length) return null;
			const now = Date.now();
			let current: number | null = null;
			for (const e of ts) {
				if (new Date(e.ts).getTime() <= now) {
					current = e.val;
				} else {
					break;
				}
			}
			return current;
		},
		days(): {
			key: string;
			energy: string;
			label: string;
			align: string;
			note: string;
		}[] {
			const s = this.solar;
			if (!s) return [];
			const dayAfterTomorrow = new Date();
			dayAfterTomorrow.setDate(dayAfterTomorrow.getDate() + 2);
			const dayAfterLabel = this.weekdayLong(dayAfterTomorrow);
			const days = [
				{
					key: "today",
					data: s.today,
					label: this.$t("forecast.solar.today"),
					align: "start",
					note: this.$t("forecast.solar.remaining"),
				},
				{
					key: "tomorrow",
					data: s.tomorrow,
					label: this.$t("forecast.solar.tomorrow"),
					align: "center",
					note: "",
				},
				{
					key: "dayAfterTomorrow",
					data: s.dayAfterTomorrow,
					label: dayAfterLabel,
					align: "end",
					note:
						s.dayAfterTomorrow && !s.dayAfterTomorrow.complete
							? this.$t("forecast.solar.partly")
							: "",
				},
			];
			return days.map((d) => ({
				key: d.key,
				energy: d.data?.energy ? this.fmtWh(d.data.energy, POWER_UNIT.AUTO) : "-",
				label: d.label,
				align: d.align,
				note: d.data?.energy ? d.note : "",
			}));
		},
	},
});
</script>

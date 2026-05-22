<template>
	<GenericModal
		:id="`loadpointSettingsModal_${id}`"
		ref="modal"
		:title="$t('main.loadpointSettings.title', [loadpoint?.title])"
		size="xl"
		data-testid="loadpoint-settings-modal"
		@open="modalVisible"
		@closed="modalInvisible"
	>
		<div class="container">
			<SmartCostLimit
				:current-limit="loadpoint?.smartCostLimit ?? null"
				:last-limit="loadpoint?.lastSmartCostLimit"
				:smart-cost-type="smartCostType"
				:currency="currency"
				is-loadpoint
				:loadpoint-id="id"
				:multiple-loadpoints="multipleLoadpoints"
				:possible="smartCostAvailable"
				:tariff="forecast?.planner"
				class="mt-2 mb-4"
			/>
			<SmartFeedInPriority
				:current-limit="loadpoint?.smartFeedInPriorityLimit ?? null"
				:last-limit="loadpoint?.lastSmartFeedInPriorityLimit"
				:currency="currency"
				:loadpoint-id="id"
				:multiple-loadpoints="multipleLoadpoints"
				:possible="smartFeedInPriorityAvailable"
				:tariff="forecast?.feedin"
				class="mt-2 mb-4"
			/>
			<LoadpointSettingsBatteryBoost
				v-if="batteryBoostAvailable"
				v-bind="batteryBoostProps"
				class="mt-2"
				@batteryboostlimit-updated="setBatteryBoostLimit"
			/>
			<h6>
				{{ $t("main.loadpointSettings.currents") }}
			</h6>
			<div v-if="phasesOptions.length" class="mb-3 row">
				<label
					:for="formId(`phases_${phasesOptions[0]}`)"
					class="col-sm-4 col-form-label pt-0"
				>
					{{ $t("main.loadpointSettings.phasesConfigured.label") }}
				</label>
				<div class="col-sm-8 pe-0">
					<p v-if="!loadpoint?.chargerPhases1p3p" class="mt-0 mb-2">
						<small>
							{{ $t("main.loadpointSettings.phasesConfigured.no1p3pSupport") }}</small
						>
					</p>
					<div v-for="phases in phasesOptions" :key="phases" class="form-check">
						<input
							:id="formId(`phases_${phases}`)"
							v-model.number="selectedPhases"
							class="form-check-input"
							type="radio"
							:name="formId('phases')"
							:value="phases"
							@change="setPhasesConfigured"
						/>
						<label class="form-check-label" :for="formId(`phases_${phases}`)">
							{{ $t(`main.loadpointSettings.phasesConfigured.phases_${phases}`) }}
							<small v-if="phases > 0">
								{{
									$t(
										`main.loadpointSettings.phasesConfigured.phases_${phases}_hint`,
										{
											min: fmtPhasePower(minCurrent, phases),
											max: fmtPhasePower(maxCurrent, phases),
										}
									)
								}}
							</small>
						</label>
					</div>
				</div>
			</div>

			<div class="mb-3 row">
				<label :for="formId('maxcurrent')" class="col-sm-4 col-form-label pt-0 pt-sm-2">
					{{ $t("main.loadpointSettings.maxCurrent.label") }}
				</label>
				<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
					<select
						:id="formId('maxcurrent')"
						v-model.number="selectedMaxCurrent"
						class="form-select form-select-sm"
						@change="setMaxCurrent"
					>
						<option
							v-for="{ value, name } in maxCurrentOptions"
							:key="value"
							:value="value"
						>
							{{ name }}
						</option>
					</select>
				</div>
			</div>

			<div class="mb-3 row">
				<label :for="formId('mincurrent')" class="col-sm-4 col-form-label pt-0 pt-sm-2">
					{{ $t("main.loadpointSettings.minCurrent.label") }}
				</label>
				<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
					<select
						:id="formId('mincurrent')"
						v-model.number="selectedMinCurrent"
						class="form-select form-select-sm"
						@change="setMinCurrent"
					>
						<option
							v-for="{ value, name } in minCurrentOptions"
							:key="value"
							:value="value"
						>
							{{ name }}
						</option>
					</select>
				</div>
			</div>

			<div v-if="phaseSwitchingPossible" class="mb-3 row">
				<label :for="formId('perphase')" class="col-sm-4 col-form-label pt-0 pt-sm-2">
					{{ $t("main.loadpointSettings.perPhase.label") }}
				</label>
				<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
					<div class="form-check form-switch">
						<input
							:id="formId('perphase')"
							v-model="showPerPhase"
							class="form-check-input"
							type="checkbox"
							role="switch"
							@change="togglePerPhase"
						/>
					</div>
				</div>
			</div>

			<template v-if="phaseSwitchingPossible && showPerPhase">
				<div class="mb-3 row align-items-center">
					<label
						:for="formId('mincurrent1p')"
						class="col-sm-4 col-form-label pt-0 pt-sm-2 ps-4"
					>
						{{ $t("main.loadpointSettings.perPhase.min1p") }}
					</label>
					<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
						<select
							:id="formId('mincurrent1p')"
							v-model.number="selectedMinCurrent1p"
							class="form-select form-select-sm"
							@change="setPerPhase('mincurrent1p', selectedMinCurrent1p)"
						>
							<option :value="0">
								{{ $t("main.loadpointSettings.perPhase.default") }}
							</option>
							<option
								v-for="{ value, name } in perPhaseOptions"
								:key="value"
								:value="value"
							>
								{{ name }}
							</option>
						</select>
					</div>
				</div>

				<div class="mb-3 row align-items-center">
					<label
						:for="formId('maxcurrent1p')"
						class="col-sm-4 col-form-label pt-0 pt-sm-2 ps-4"
					>
						{{ $t("main.loadpointSettings.perPhase.max1p") }}
					</label>
					<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
						<select
							:id="formId('maxcurrent1p')"
							v-model.number="selectedMaxCurrent1p"
							class="form-select form-select-sm"
							@change="setPerPhase('maxcurrent1p', selectedMaxCurrent1p)"
						>
							<option :value="0">
								{{ $t("main.loadpointSettings.perPhase.default") }}
							</option>
							<option
								v-for="{ value, name } in perPhaseOptions"
								:key="value"
								:value="value"
							>
								{{ name }}
							</option>
						</select>
					</div>
				</div>

				<div class="mb-3 row align-items-center">
					<label
						:for="formId('mincurrent3p')"
						class="col-sm-4 col-form-label pt-0 pt-sm-2 ps-4"
					>
						{{ $t("main.loadpointSettings.perPhase.min3p") }}
					</label>
					<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
						<select
							:id="formId('mincurrent3p')"
							v-model.number="selectedMinCurrent3p"
							class="form-select form-select-sm"
							@change="setPerPhase('mincurrent3p', selectedMinCurrent3p)"
						>
							<option :value="0">
								{{ $t("main.loadpointSettings.perPhase.default") }}
							</option>
							<option
								v-for="{ value, name } in perPhaseOptions"
								:key="value"
								:value="value"
							>
								{{ name }}
							</option>
						</select>
					</div>
				</div>

				<div class="mb-3 row align-items-center">
					<label
						:for="formId('maxcurrent3p')"
						class="col-sm-4 col-form-label pt-0 pt-sm-2 ps-4"
					>
						{{ $t("main.loadpointSettings.perPhase.max3p") }}
					</label>
					<div class="col-sm-8 col-lg-4 pe-0 d-flex align-items-center">
						<select
							:id="formId('maxcurrent3p')"
							v-model.number="selectedMaxCurrent3p"
							class="form-select form-select-sm"
							@change="setPerPhase('maxcurrent3p', selectedMaxCurrent3p)"
						>
							<option :value="0">
								{{ $t("main.loadpointSettings.perPhase.default") }}
							</option>
							<option
								v-for="{ value, name } in perPhaseOptions"
								:key="value"
								:value="value"
							>
								{{ name }}
							</option>
						</select>
					</div>
				</div>
			</template>
		</div>
	</GenericModal>
</template>

<script lang="ts">
import collector from "@/mixins/collector.ts";
import formatter from "@/mixins/formatter";
import GenericModal from "../Helper/GenericModal.vue";
import SmartCostLimit from "../Tariff/SmartCostLimit.vue";
import SmartFeedInPriority from "../Tariff/SmartFeedInPriority.vue";
import SettingsBatteryBoost from "./SettingsBatteryBoost.vue";
import { defineComponent, type PropType } from "vue";
import { PHASES, CURRENCY, SMART_COST_TYPE, type Forecast, type UiLoadpoint } from "@/types/evcc";
import api from "@/api";

const V = 230;

const range = (start: number, stop: number, step = -1) =>
	Array.from({ length: (stop - start) / step + 1 }, (_, i) => start + i * step);

const insertSorted = (arr: number[], num: number) => {
	const uniqueSet = new Set(arr);
	uniqueSet.add(num);
	return [...uniqueSet].sort((a, b) => b - a);
};

// TODO: add max physical current to loadpoint (config ui) and only allow user to select values in side that range (main ui, here)
const MAX_CURRENT = 64;

const { AUTO, THREE_PHASES, ONE_PHASE } = PHASES;

export default defineComponent({
	name: "LoadpointSettingsModal",
	components: {
		GenericModal,
		SmartCostLimit,
		SmartFeedInPriority,
		LoadpointSettingsBatteryBoost: SettingsBatteryBoost,
	},
	mixins: [formatter, collector],
	props: {
		loadpoints: { type: Array as PropType<UiLoadpoint[]>, default: () => [] },
		batteryConfigured: Boolean,
		smartCostType: String as PropType<SMART_COST_TYPE>,
		smartCostAvailable: Boolean,
		smartFeedInPriorityAvailable: Boolean,
		tariffGrid: Number,
		currency: String as PropType<CURRENCY>,
		multipleLoadpoints: Boolean,
		forecast: Object as PropType<Forecast>,
	},
	data() {
		return {
			id: undefined as string | undefined,
			selectedMaxCurrent: undefined as number | undefined,
			selectedMinCurrent: undefined as number | undefined,
			selectedPhases: undefined as number | undefined,
			isModalVisible: false,
			showPerPhase: false,
			selectedMinCurrent1p: 0,
			selectedMaxCurrent1p: 0,
			selectedMinCurrent3p: 0,
			selectedMaxCurrent3p: 0,
		};
	},
	computed: {
		loadpoint() {
			return this.loadpoints.find((loadpoint) => loadpoint.id === this.id);
		},
		maxCurrent() {
			return this.loadpoint?.maxCurrent;
		},
		minCurrent() {
			return this.loadpoint?.minCurrent;
		},
		minCurrent1p() {
			return this.loadpoint?.minCurrent1p ?? null;
		},
		maxCurrent1p() {
			return this.loadpoint?.maxCurrent1p ?? null;
		},
		minCurrent3p() {
			return this.loadpoint?.minCurrent3p ?? null;
		},
		maxCurrent3p() {
			return this.loadpoint?.maxCurrent3p ?? null;
		},
		// per-phase limits only make sense on chargers that can switch 1p<->3p
		phaseSwitchingPossible() {
			return !!this.loadpoint?.chargerPhases1p3p;
		},
		// full current range, used for all four per-phase selects
		perPhaseOptions() {
			return range(MAX_CURRENT, 1).map((value) => ({
				value,
				name: `${this.fmtNumber(value, undefined)} A`,
			}));
		},
		batteryBoostLimit() {
			return this.loadpoint?.batteryBoostLimit;
		},
		phasesConfigured() {
			return this.loadpoint?.phasesConfigured;
		},
		phasesOptions() {
			if (this.loadpoint?.chargerSinglePhase) {
				return [];
			}
			if (this.loadpoint?.chargerPhases1p3p) {
				// automatic switching
				return [AUTO, THREE_PHASES, ONE_PHASE];
			}
			// 1p or 3p possible
			return [THREE_PHASES, ONE_PHASE];
		},
		batteryBoostProps() {
			return this.collectProps(SettingsBatteryBoost);
		},
		maxPhases() {
			if (this.loadpoint?.chargerPhases1p3p && this.phasesConfigured === AUTO) {
				return THREE_PHASES;
			}
			return this.phasesConfigured;
		},
		minPhases() {
			if (this.loadpoint?.chargerPhases1p3p && this.phasesConfigured === AUTO) {
				return ONE_PHASE;
			}
			return this.phasesConfigured;
		},
		minCurrentOptions() {
			const opt1 = [...range(Math.floor(this.maxCurrent ?? 0), 1), 0.5, 0.25, 0.125];
			// ensure that current value is always included
			const opt2 = insertSorted(opt1, this.minCurrent ?? 0);
			return opt2.map((value) => this.currentOption(value, value === 6, this.minPhases));
		},
		maxCurrentOptions() {
			const opt1 = range(MAX_CURRENT, Math.ceil(this.minCurrent ?? 0));
			// ensure that current value is always included
			const opt2 = insertSorted(opt1, this.maxCurrent ?? 0);
			return opt2.map((value) => this.currentOption(value, value === 16, this.maxPhases));
		},
		batteryBoostAvailable() {
			return this.batteryConfigured;
		},
	},
	watch: {
		maxCurrent(value) {
			this.selectedMaxCurrent = value;
		},
		minCurrent(value) {
			this.selectedMinCurrent = value;
		},
		phasesConfigured(value) {
			this.selectedPhases = value;
		},
		minCurrent1p(value) {
			this.selectedMinCurrent1p = value ?? 0;
		},
		maxCurrent1p(value) {
			this.selectedMaxCurrent1p = value ?? 0;
		},
		minCurrent3p(value) {
			this.selectedMinCurrent3p = value ?? 0;
		},
		maxCurrent3p(value) {
			this.selectedMaxCurrent3p = value ?? 0;
		},
	},
	methods: {
		open(loadpointId: string) {
			this.id = loadpointId;
			this.selectedPhases = this.phasesConfigured;
			this.selectedMaxCurrent = this.maxCurrent;
			this.selectedMinCurrent = this.minCurrent;
			this.selectedMinCurrent1p = this.minCurrent1p ?? 0;
			this.selectedMaxCurrent1p = this.maxCurrent1p ?? 0;
			this.selectedMinCurrent3p = this.minCurrent3p ?? 0;
			this.selectedMaxCurrent3p = this.maxCurrent3p ?? 0;
			// expand the per-phase section if any override is already set
			this.showPerPhase =
				!!this.minCurrent1p ||
				!!this.maxCurrent1p ||
				!!this.minCurrent3p ||
				!!this.maxCurrent3p;
			const modalRef = this.$refs["modal"] as InstanceType<typeof GenericModal> | undefined;
			modalRef?.open();
		},
		apiPath(func: string) {
			return "loadpoints/" + this.id + "/" + func;
		},
		fmtPhasePower(current?: number, phases?: PHASES) {
			return this.fmtW(V * (current || 0) * (phases || 0));
		},
		formId(name: string) {
			return `loadpoint_${this.id}_${name}`;
		},
		setMaxCurrent() {
			api.post(this.apiPath("maxcurrent") + "/" + this.selectedMaxCurrent);
		},
		setMinCurrent() {
			api.post(this.apiPath("mincurrent") + "/" + this.selectedMinCurrent);
		},
		setPhasesConfigured() {
			api.post(this.apiPath("phases") + "/" + this.selectedPhases);
		},
		// setPerPhase posts a per-phase override, or DELETEs it when set to the
		// "default" sentinel (0) so the loadpoint falls back to its global min/max
		setPerPhase(func: string, value: number) {
			if (value > 0) {
				api.post(this.apiPath(func) + "/" + value);
			} else {
				api.delete(this.apiPath(func));
			}
		},
		// togglePerPhase clears all four overrides when the section is switched off
		togglePerPhase() {
			if (!this.showPerPhase) {
				this.selectedMinCurrent1p = 0;
				this.selectedMaxCurrent1p = 0;
				this.selectedMinCurrent3p = 0;
				this.selectedMaxCurrent3p = 0;
				["mincurrent1p", "maxcurrent1p", "mincurrent3p", "maxcurrent3p"].forEach((func) =>
					api.delete(this.apiPath(func))
				);
			}
		},
		setBatteryBoostLimit(limit: number) {
			api.post(this.apiPath("batteryboostlimit") + "/" + limit);
		},
		currentOption(current: number, isDefault: boolean, phases?: number) {
			const kw = this.fmtPhasePower(current, phases);
			let name = `${this.fmtNumber(current, undefined)} A (${kw})`;
			if (isDefault) {
				name += ` [${this.$t("main.loadpointSettings.default")}]`;
			}
			return { value: current, name };
		},
		modalVisible() {
			this.isModalVisible = true;
		},
		modalInvisible() {
			this.isModalVisible = false;
		},
	},
});
</script>
<style scoped>
.container {
	margin-left: calc(var(--bs-gutter-x) * -0.5);
	margin-right: calc(var(--bs-gutter-x) * -0.5);
}

.container h4:first-child {
	margin-top: 0 !important;
}

.custom-select-inline {
	display: inline-block !important;
}
</style>

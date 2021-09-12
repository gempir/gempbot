import { Outcome, RawPredictionLog } from "../hooks/usePredictionLogs";

export class PredictionLog {
    constructor(
        public readonly ID: string,
        public readonly OwnerTwitchID: string,
        public readonly Title: string,
        public readonly WinningOutcomeID: string,
        public readonly Status: string,
        public readonly StartedAt: Date,
        public readonly LockedAt: Date,
        public readonly EndedAt: Date,
        public readonly Outcomes: Outcome[]
    ) { }

    public static fromObject(data: RawPredictionLog): PredictionLog {
        return new PredictionLog(data.ID, data.OwnerTwitchID, data.Title, data.WinningOutcomeID, data.Status, new Date(data.StartedAt), new Date(data.LockedAt), new Date(data.EndedAt), data.Outcomes);
    }

    public getWinningOutcome(): Outcome | undefined {
        let winner;

        for (const outcome of this.Outcomes) {
            if (this.WinningOutcomeID === outcome.ID) {
                winner = outcome;
            }
        }

        return winner
    }
}
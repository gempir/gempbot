import React, { useState } from "react";
import styled from "styled-components";
import { useRewards } from "../hooks/useRewards";
import { useUserConfig } from "../hooks/useUserConfig";
import { Menu } from "./Menu";
import { Rewards } from "./Rewards";

export function Dashboard() {
    const [classNames, setClassNames] = useState(["userConfig"]);

    const [userCfg, setUserConfig] = useUserConfig(() => {
        const newClassNames = classNames.slice();
        newClassNames.push("saved");
        setClassNames(newClassNames)

        setTimeout(() => {
            setClassNames(["userConfig"]);
        }, 500);
    }, () => {
        const newClassNames = classNames.slice();
        newClassNames.push("error");
        setClassNames(newClassNames)

        setTimeout(() => {
            setClassNames(["userConfig"]);
        }, 500);
    });

    const [rewards] = useRewards();

    return <DashboardContainer>
        {userCfg && <Menu userConfig={userCfg} setUserConfig={setUserConfig} />}
        {<div className={classNames.join(" ")}>
            {userCfg && <div className={"redemption"}>
                <img src="/images/bttv.png" alt={"bttv"} />
                <label className="switch">
                    <input type="checkbox" checked={userCfg.Redemptions.Bttv.Active} onChange={(e) => {
                        const newConfig = JSON.parse(JSON.stringify(userCfg));
                        newConfig.Redemptions.Bttv.Active = e.target.checked;

                        setUserConfig(newConfig);
                    }} />
                    <span className="slider round"></span>
                </label>
                <div className="redemption-title">
                    <span>Channel Points Reward Name</span>
                    <input type="text" value={userCfg.Redemptions.Bttv.Title} spellCheck={false} onChange={(e) => {
                        const newConfig = JSON.parse(JSON.stringify(userCfg));
                        newConfig.Redemptions.Bttv.Title = e.target.value;

                        setUserConfig(newConfig);
                    }} />
                </div>
                <span className="hint">
                    make sure <strong>gempbot</strong> is bttv editor
                </span>
            </div>}
        </div>}
        <Rewards rewards={rewards} />
    </DashboardContainer>
}

const DashboardContainer = styled.div`
    margin-top: 5rem;
    margin-left: 1rem;
    margin-right: 1rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    grid-gap: 2rem;

    .userConfig {
        padding-bottom: 2rem;
        background: var(--bg);
        transition: background-color ease-in-out 0.2s;

        &.saved {
            background: var(--theme);
        }

        &.error {
            background: var(--danger);
        }
    }

    .redemption {
        display: flex;
        align-items: center;
        background: var(--bg-bright);
        border: 1px solid var(--bg-brighter);
        padding: 0.5rem;

        img {
            max-height: 3rem;
            margin-left: 1rem;
            margin-right: 2rem;
        }

        .redemption-title {
            position: relative;

            span {
                position: absolute;
                top: -13px;
                left: 18px;
                font-size: 11px;
            }

            input {
                margin: 0;
                padding: 0;
                margin-left: 1rem;
                font-size: 1rem;
                background: var(--bg);
                border: 1px solid var(--bg-bright);
                padding: 5px;
                color: white;

                &:focus {
                    outline: none;
                    border: 1px solid var(--theme2);
                }
            }
        }
        
        .hint {
            margin-left: 1rem;

            strong {
                color: var(--theme-bright);
            }
        }
    }

    /* The switch - the box around the slider */
    .switch {
        position: relative;
        display: inline-block;
        width: 60px;
        height: 34px;

        input {
            opacity: 0;
            width: 0;
            height: 0;
        }
    }
    /* The slider */
    .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: #ccc;
        -webkit-transition: .4s;
        transition: .4s;

        &:before {
            position: absolute;
            content: "";
            height: 26px;
            width: 26px;
            left: 4px;
            bottom: 4px;
            background-color: white;
            -webkit-transition: .4s;
            transition: .4s;
        }
    }

    input:checked + .slider {
        background-color: var(--theme-bright);
    }

    input:focus + .slider {
     box-shadow: 0 0 1px var(--theme-bright);
    }

    input:checked + .slider:before {
        -webkit-transform: translateX(26px);
        -ms-transform: translateX(26px);
        transform: translateX(26px);
    }

    .slider.round {
        border-radius: 34px;

        &:before {
            border-radius: 50%;
        }
    }
`;
import { useState } from "react";
import styled from "styled-components";
import { useUserConfig } from "../hooks/useUserConfig";

export function Dashboard() {
    const [classNames, setClassNames] = useState(["userConfig"]);

    const [userCfg, setConfig] = useUserConfig(() => {
        const newClassNames = classNames.slice();
        newClassNames.push("saved");
        setClassNames(newClassNames)

        setTimeout(() => {
            setClassNames(["userConfig"]);
        }, 500);
    });

    return <DashboardContainer>
        {userCfg && <div className={classNames.join(" ")}>
            <div className={"redemption"}>
                <img src="/images/bttv.png" alt={"bttv"} />
                <label className="switch">
                    <input type="checkbox" checked={userCfg.Redemptions.Bttv.Active} onChange={(e) => {
                        const newConfig = JSON.parse(JSON.stringify(userCfg));
                        newConfig.Redemptions.Bttv.Active = e.target.checked;

                        setConfig(newConfig);
                    }} />
                    <span className="slider round"></span>
                </label>
                <div className="redemption-title">
                    <span>Channel Points Reward Name</span>
                    <input type="text" value={userCfg.Redemptions.Bttv.Title} spellCheck={false} onChange={(e) => {
                        const newConfig = JSON.parse(JSON.stringify(userCfg));
                        newConfig.Redemptions.Bttv.Title = e.target.value;

                        setConfig(newConfig);
                    }} />
                </div>
                <span className="hint">
                    make sure <strong>gempbot</strong> is bttv editor
                </span>
            </div>
        </div>}
    </DashboardContainer>
}

const DashboardContainer = styled.div`
    margin-top: 5rem;
    margin-left: 1rem;
    margin-right: 1rem;

    .userConfig {
        padding-bottom: 2rem;
        background: var(--bg);
        transition: background-color ease-in-out 0.2s;

        &.saved {
            background: var(--theme);
        }
    }

    .redemption {
        display: flex;
        align-items: center;
        background: var(--bg-dark);
        padding: 1rem;

        img {
            max-height: 3rem;
            margin-left: 1rem;
            margin-right: 2rem;
        }

        .redemption-title {
            position: relative;

            span {
                position: absolute;
                top: -15px;
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
    }

    /* Hide default HTML checkbox */
    .switch input {
        opacity: 0;
        width: 0;
        height: 0;
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
    }

    .slider:before {
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

    input:checked + .slider {
        background-color: #2196F3;
    }

    input:focus + .slider {
     box-shadow: 0 0 1px #2196F3;
    }

    input:checked + .slider:before {
        -webkit-transform: translateX(26px);
        -ms-transform: translateX(26px);
        transform: translateX(26px);
    }

    /* Rounded sliders */
    .slider.round {
        border-radius: 34px;
    }

    .slider.round:before {
        border-radius: 50%;
    }
`;
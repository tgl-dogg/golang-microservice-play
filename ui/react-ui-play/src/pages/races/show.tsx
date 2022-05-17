import { useShow } from "@pankod/refine-core";
import { Show, Typography, Tag } from "@pankod/refine-antd";

//import { IRace } from "interfaces";

const { Title, Text } = Typography;

export const RaceShow = () => {
    const { queryResult } = useShow();
    const { data, isLoading } = queryResult;
    const record = data?.data;

    return (
        <Show isLoading={isLoading}>
            <Title level={5}>Name</Title>
            <Text>{record?.name}</Text>

            <Title level={5}>Description</Title>
            <Text><Tag>{record?.description}</Tag></Text>
        </Show>
    );
};